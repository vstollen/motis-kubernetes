/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	motisv1alpha1 "github.com/vstollen/motis-operator/api/v1alpha1"
)

// DatasetReconciler reconciles a Dataset object
type DatasetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=motis.motis-project.de,resources=datasets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=motis.motis-project.de,resources=datasets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=motis.motis-project.de,resources=datasets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *DatasetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	dataset := &motisv1alpha1.Dataset{}
	if err := r.Get(ctx, req.NamespacedName, dataset); err != nil {
		if errors.IsNotFound(err) {
			log.Info("dataset resource not found. Ignoring since object was likely deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Dataset")
		return ctrl.Result{}, err
	}

	if dataset.Status.InputVolume == nil {
		log.Info("No input volume claimed. Creating PVC")
		if err := r.createInputPVC(ctx, dataset, log); err != nil {
			log.Error(err, "unable to creat input PVC")
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, req.NamespacedName, dataset); err != nil {
			log.Error(err, "Failed to re-fetch dataset")
			return ctrl.Result{}, err
		}
	}

	if dataset.Status.DataVolume == nil {
		log.Info("No data volume claimed. Creating PVC")
		return r.createDataPVC(ctx, dataset, log)
	}

	if !dataset.Status.StartedProcessing {
		result, err := r.createProcessingJob(ctx, dataset, log)
		if err != nil {
			return result, err
		}

		dataset.Status.StartedProcessing = true
		err = r.Client.Status().Update(ctx, dataset)
		if err != nil {
			log.Error(err, "unable to update status with started processing job")
			return ctrl.Result{}, err
		}
	}

	if !dataset.Status.FinishedProcessing {
		var processingJob batchv1.Job
		if err := r.Get(ctx, client.ObjectKey{
			Namespace: req.Namespace,
			Name:      dataset.Status.ProcessingJobName,
		}, &processingJob); err != nil {
			log.Error(err, "unable to list child jobs")
			return ctrl.Result{}, err
		}

		for _, condition := range processingJob.Status.Conditions {
			if condition.Type == batchv1.JobComplete && condition.Status == corev1.ConditionTrue {
				dataset.Status.FinishedProcessing = true
				err := r.Client.Status().Update(ctx, dataset)
				if err != nil {
					log.Error(err, "unable to update status with finished processing job")
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{}, nil
}

func (r *DatasetReconciler) createProcessingJob(ctx context.Context, dataset *motisv1alpha1.Dataset, log logr.Logger) (ctrl.Result, error) {
	processJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "motis-processing-job-",
			Namespace:    dataset.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "motis-init",
							Image: "ghcr.io/vstollen/motis-init:0.1.1",
							//Image:   "busybox",
							//Command: []string{"tail"},
							//Args:    []string{"-f", "/dev/null"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/config",
								},
								{
									Name:      "input-volume",
									MountPath: "/input",
								},
							},
						},
					},
					Containers: []corev1.Container{
						{
							Name:    "motis",
							Image:   "ghcr.io/motis-project/motis:latest",
							Command: []string{"/motis/motis", "--system_config", "/system_config.ini", "-c", "/config/config.ini", "--mode", "test"},
							VolumeMounts: []corev1.VolumeMount{{
								Name:      "data-volume",
								MountPath: "/data",
							},
								{
									Name:      "input-volume",
									MountPath: "/input",
								},
								{
									Name:      "config",
									MountPath: "/config",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "data-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: dataset.Status.DataVolume,
							},
						},
						{
							Name: "input-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: dataset.Status.InputVolume,
							},
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: dataset.Spec.Config,
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}

	err := ctrl.SetControllerReference(dataset, processJob, r.Scheme)
	if err != nil {
		log.Error(err, "unable to set controller reference on data pvc")
		return ctrl.Result{}, err
	}

	err = r.Client.Create(ctx, processJob)
	if err != nil {
		log.Error(err, "unable to create processing job")
		return ctrl.Result{}, err
	}

	dataset.Status.ProcessingJobName = processJob.Name
	err = r.Client.Status().Update(ctx, dataset)
	if err != nil {
		log.Error(err, "unable to update status with processing job name")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, err
}

func (r *DatasetReconciler) createInputPVC(ctx context.Context, dataset *motisv1alpha1.Dataset, log logr.Logger) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "motis-input-claim-",
			Namespace:    dataset.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse("10Gi"),
				},
			},
		},
	}

	err := ctrl.SetControllerReference(dataset, pvc, r.Scheme)
	if err != nil {
		log.Error(err, "unable to set controller reference on input pvc")
		return err
	}

	err = r.Client.Create(ctx, pvc)
	if err != nil {
		log.Error(err, "unable to create input volume pvc")
		return err
	}

	dataset.Status.InputVolume = &corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: pvc.Name,
	}

	err = r.Client.Status().Update(ctx, dataset)
	if err != nil {
		log.Error(err, "unable to update status with new input pvc")
		return err
	}
	return nil
}

func (r *DatasetReconciler) createDataPVC(ctx context.Context, dataset *motisv1alpha1.Dataset, log logr.Logger) (ctrl.Result, error) {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "motis-data-claim-",
			Namespace:    dataset.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse("10Gi"),
				},
			},
		},
	}

	err := ctrl.SetControllerReference(dataset, pvc, r.Scheme)
	if err != nil {
		log.Error(err, "unable to set controller reference on data pvc")
		return ctrl.Result{}, err
	}

	err = r.Client.Create(ctx, pvc)
	if err != nil {
		log.Error(err, "unable to create data volume pvc")
		return ctrl.Result{}, err
	}

	dataset.Status.DataVolume = &corev1.PersistentVolumeClaimVolumeSource{
		ClaimName: pvc.Name,
	}

	err = r.Client.Status().Update(ctx, dataset)
	if err != nil {
		log.Error(err, "unable to update status with new data pvc")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatasetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&motisv1alpha1.Dataset{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
