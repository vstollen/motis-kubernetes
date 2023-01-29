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
	"k8s.io/apimachinery/pkg/types"

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
	log.Info("Starting reconciliation")

	dataset := &motisv1alpha1.Dataset{}
	log.Info("Fetching dataset")
	if err := r.Get(ctx, req.NamespacedName, dataset); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Dataset resource not found. Ignoring since object was likely deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Error fetching dataset")
		return ctrl.Result{}, err
	}

	inputVolume := &corev1.PersistentVolumeClaim{}
	log.Info("Fetching input volume")
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name + "-input", Namespace: req.Namespace}, inputVolume); client.IgnoreNotFound(err) != nil {
		log.Error(err, "Error retrieving input volume")
		return ctrl.Result{}, err
	}

	dataVolume := &corev1.PersistentVolumeClaim{}
	log.Info("Fetching data volume")
	if err := r.Get(ctx, types.NamespacedName{Name: req.Name + "-data", Namespace: req.Namespace}, dataVolume); client.IgnoreNotFound(err) != nil {
		log.Error(err, "Error retrieving data volume")
		return ctrl.Result{}, err
	}

	processingJob := &batchv1.Job{}
	log.Info("Fetching processing job")
	if err := r.Get(ctx, types.NamespacedName{Name: dataset.Name, Namespace: dataset.Namespace}, processingJob); client.IgnoreNotFound(err) != nil {
		log.Error(err, "Error retrieving processing job")
		return ctrl.Result{}, err
	}

	log.Info("Updating status")
	if err := r.updateStatus(ctx, req, dataset, inputVolume, dataVolume, processingJob, log); err != nil {
		log.Error(err, "Error updating status")
	}

	if inputVolume == nil || inputVolume.UID == "" {
		log.Info("No input volume claimed. Creating PVC")
		if err := r.createInputPVC(ctx, dataset, log); err != nil {
			log.Error(err, "unable to create input PVC")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if dataVolume == nil || dataVolume.UID == "" {
		log.Info("No data volume claimed. Creating PVC")
		err := r.createDataPVC(ctx, dataset, log)
		if err != nil {
			log.Error(err, "unable to create data PVC")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if processingJob == nil || processingJob.UID == "" {
		log.Info("No processing job found. Creating new processing job")
		if err := r.createProcessingJob(ctx, dataset, log); err != nil {
			log.Error(err, "Error creating processing job")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *DatasetReconciler) updateStatus(ctx context.Context, req ctrl.Request, dataset *motisv1alpha1.Dataset, inputVolume *corev1.PersistentVolumeClaim, dataVolume *corev1.PersistentVolumeClaim, processingJob *batchv1.Job, log logr.Logger) error {
	if inputVolume != nil {
		dataset.Status.InputVolume = &corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: inputVolume.Name},
		}
	} else {
		dataset.Status.InputVolume = nil
	}

	if dataVolume != nil {
		dataset.Status.DataVolume = &corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: dataVolume.Name},
		}
	} else {
		dataset.Status.DataVolume = nil
	}

	readyCondition := &motisv1alpha1.DatasetCondition{
		Type:   motisv1alpha1.DatasetReady,
		Status: corev1.ConditionUnknown,
	}

	if processingJob != nil {
		for _, condition := range processingJob.Status.Conditions {
			if condition.Type == batchv1.JobComplete {
				switch condition.Status {
				case corev1.ConditionTrue:
					readyCondition.Status = corev1.ConditionTrue
					break
				case corev1.ConditionFalse:
					readyCondition.Status = corev1.ConditionFalse
					break
				}
			}
		}
	}

	dataset.Status.Conditions = []motisv1alpha1.DatasetCondition{*readyCondition}

	if err := r.Client.Status().Update(ctx, dataset); err != nil {
		log.Error(err, "Error updating status")
		return err
	}

	if err := r.Get(ctx, req.NamespacedName, dataset); err != nil {
		log.Error(err, "Error fetching updated dataset")
		return err
	}
	return nil
}

func (r *DatasetReconciler) createInputPVC(ctx context.Context, dataset *motisv1alpha1.Dataset, log logr.Logger) error {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset.Name + "-input",
			Namespace: dataset.Namespace,
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

	return nil
}

func (r *DatasetReconciler) createDataPVC(ctx context.Context, dataset *motisv1alpha1.Dataset, log logr.Logger) error {
	pvc := r.dataPvcForDataset(dataset)

	if err := ctrl.SetControllerReference(dataset, pvc, r.Scheme); err != nil {
		log.Error(err, "unable to set controller reference on data pvc")
		return err
	}

	if err := r.Client.Create(ctx, pvc); err != nil {
		log.Error(err, "unable to create data volume pvc")
		return err
	}

	return nil
}

func (r *DatasetReconciler) createProcessingJob(ctx context.Context, dataset *motisv1alpha1.Dataset, log logr.Logger) error {
	job := r.processingJobForDataset(dataset)

	if err := ctrl.SetControllerReference(dataset, job, r.Scheme); err != nil {
		log.Error(err, "unable to set controller reference on processing job")
		return err
	}

	if err := r.Client.Create(ctx, job); err != nil {
		log.Error(err, "unable to create processing job")
		return err
	}

	return nil
}

func (r *DatasetReconciler) dataPvcForDataset(dataset *motisv1alpha1.Dataset) *corev1.PersistentVolumeClaim {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset.Name + "-data",
			Namespace: dataset.Namespace,
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

	return pvc
}

func (r *DatasetReconciler) processingJobForDataset(dataset *motisv1alpha1.Dataset) *batchv1.Job {
	processJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dataset.Name,
			Namespace: dataset.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					InitContainers: []corev1.Container{
						{
							Name:  "motis-init",
							Image: "ghcr.io/vstollen/motis-init:0.1.1",
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
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: dataset.Name + "-data"},
							},
						},
						{
							Name: "input-volume",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: dataset.Name + "-input"},
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
	return processJob
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatasetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&motisv1alpha1.Dataset{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
