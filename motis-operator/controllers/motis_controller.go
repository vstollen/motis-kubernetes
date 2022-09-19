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
	"errors"
	corev1 "k8s.io/api/core/v1"
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
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Motis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *MotisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	motis := &motisv1alpha1.Motis{}
	if err := r.Get(ctx, req.NamespacedName, motis); err != nil {
		log.Error(err, "unable to fetch Motis")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if motis.Spec.InputVolume == nil {
		log.Error(errors.New("InputVolume not set"), "InputVolume must be specified")
	}

	if motis.Status.DatasetName == "" {
		dataset := &motisv1alpha1.Dataset{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "motis-dataset-",
				Namespace:    req.Namespace,
			},
			Spec: motisv1alpha1.DatasetSpec{
				InputVolume: motis.Spec.InputVolume,
				Config:      motis.Spec.Config,
			},
		}

		err := ctrl.SetControllerReference(motis, dataset, r.Scheme)
		if err != nil {
			log.Error(err, "unable to set controller reference on dataset")
			return ctrl.Result{}, err
		}

		err = r.Client.Create(ctx, dataset)
		if err != nil {
			log.Error(err, "unable to create motis dataset")
			return ctrl.Result{}, err
		}

		motis.Status.DatasetName = dataset.Name
		err = r.Client.Status().Update(ctx, motis)
		if err != nil {
			log.Error(err, "unable to update status with dataset name")
			return ctrl.Result{}, err
		}
	}

	dataset := &motisv1alpha1.Dataset{}
	if err := r.Get(ctx, client.ObjectKey{
		Name:      motis.Status.DatasetName,
		Namespace: req.Namespace,
	}, dataset); err != nil {
		log.Error(err, "unable to fetch dataset")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if motis.Status.PodName == "" && dataset.Status.FinishedProcessing {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "motis-",
				Namespace:    req.Namespace,
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
						Name: "data-volume",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: dataset.Status.DataVolume,
						},
					},
					{
						Name: "input-volume",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: motis.Spec.InputVolume,
						},
					},
					{
						Name: "config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: motis.Spec.Config,
						},
					},
				},
			},
		}

		err := ctrl.SetControllerReference(motis, pod, r.Scheme)
		if err != nil {
			log.Error(err, "unable to set controller reference on motis pod")
			return ctrl.Result{}, err
		}

		err = r.Client.Create(ctx, pod)
		if err != nil {
			log.Error(err, "unable to create motis pod")
			return ctrl.Result{}, err
		}

		motis.Status.PodName = pod.Name
		err = r.Client.Status().Update(ctx, motis)
		if err != nil {
			log.Error(err, "unable to update status with motis pod name")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MotisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&motisv1alpha1.Motis{}).
		Owns(&motisv1alpha1.Dataset{}).
		Complete(r)
}