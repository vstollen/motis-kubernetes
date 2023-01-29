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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	motisv1alpha1 "github.com/vstollen/motis-operator/api/v1alpha1"
)

// MotisReconciler reconciles a Motis object
type MotisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=motis.motis-project.de,resources=motis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=motis.motis-project.de,resources=motis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=motis.motis-project.de,resources=motis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *MotisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	motis := &motisv1alpha1.Motis{}
	if err := r.Get(ctx, req.NamespacedName, motis); err != nil {
		if errors.IsNotFound(err) {
			log.Info("motis resource not found. Ignoring since object was likely deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get motis resource")
		return ctrl.Result{}, err
	}

	dataset := &motisv1alpha1.Dataset{}
	if err := r.Get(ctx, req.NamespacedName, dataset); client.IgnoreNotFound(err) != nil {
		log.Error(err, "Failed to get Dataset")
		return ctrl.Result{}, err
	}

	if dataset == nil || dataset.UID == "" {
		if err := r.createDataset(ctx, motis, log); err != nil {
			log.Error(err, "Failed to create new Dataset")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, req.NamespacedName, deployment); client.IgnoreNotFound(err) != nil {
		log.Error(err, "Failed to get Motis deployment")
		return ctrl.Result{}, err
	}

	if dataset.HasFinishedProcessing() && (deployment == nil || deployment.UID == "") {
		log.Info("Dataset has finished processing and there currently is no deployment. Starting new deployment.")
		if err := r.createDeployment(ctx, motis, dataset, log); err != nil {
			log.Error(err, "Error creating Motis deployment")
		}
	}

	return ctrl.Result{}, nil
}

func (r *MotisReconciler) createDataset(ctx context.Context, motis *motisv1alpha1.Motis, log logr.Logger) error {
	dataset := datasetForMotis(motis)

	if err := ctrl.SetControllerReference(motis, dataset, r.Scheme); err != nil {
		return err
	}

	log.Info("Creating new Dataset", "Dataset.Namespace", dataset.Namespace, "Dataset.Name", dataset.Name)
	if err := r.Create(ctx, dataset); err != nil {
		log.Error(err, "Failed to create new Dataset", "Dataset.Namespace", dataset.Namespace, "Dataset.Name", dataset.Name)
		return err
	}

	return nil
}

func (r *MotisReconciler) createDeployment(ctx context.Context, motis *motisv1alpha1.Motis, dataset *motisv1alpha1.Dataset, log logr.Logger) error {
	pod := deploymentForMotis(motis, dataset)

	if err := ctrl.SetControllerReference(motis, pod, r.Scheme); err != nil {
		log.Error(err, "Error setting controller reference for deployment")
		return err
	}

	log.Info("Creating new motis deployment")
	if err := r.Create(ctx, pod); err != nil {
		log.Error(err, "Failed to create motis deployment")
		return err
	}

	return nil
}

func datasetForMotis(motis *motisv1alpha1.Motis) *motisv1alpha1.Dataset {
	return &motisv1alpha1.Dataset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      motis.Name,
			Namespace: motis.Namespace,
		},
		Spec: motisv1alpha1.DatasetSpec{
			Config: motis.Spec.Config,
		},
	}
}

func deploymentForMotis(motis *motisv1alpha1.Motis, dataset *motisv1alpha1.Dataset) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      motis.Name,
			Namespace: motis.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"motis-project.de/motis-deployment": motis.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"motis-project.de/motis-deployment": motis.Name,
						"motis-project.de/name":             "MotisWeb",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "motis",
							Image:   "ghcr.io/motis-project/motis:latest",
							Command: []string{"/motis/motis", "--system_config", "/system_config.ini", "-c", "/config/config.ini"},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
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
							Name:         "data-volume",
							VolumeSource: *dataset.Status.DataVolume,
						},
						{
							Name:         "input-volume",
							VolumeSource: *dataset.Status.InputVolume,
						},
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: motis.Spec.Config,
							},
						},
					},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *MotisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&motisv1alpha1.Motis{}).
		Owns(&motisv1alpha1.Dataset{}).
		Complete(r)
}
