/*
Copyright 2024.

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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	backupv1alpha1 "github.com/DanielKoklev/velero-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VeleroReconciler reconciles a Velero object
type VeleroReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=backup.go-learning.com,resources=veleroes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=backup.go-learning.com,resources=veleroes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=backup.go-learning.com,resources=veleroes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Velero object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *VeleroReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	velero := &backupv1alpha1.Velero{}
	err := r.Get(ctx, req.NamespacedName, velero)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Define the desired Velero deployment
	deployment := r.desiredDeployment(velero)

	// Check if the Deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, client.ObjectKey{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Update the Deployment if necessary
	if !r.deploymentEqual(found, deployment) {
		log.Info("Updating Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Update(ctx, deployment)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// Update the Velero status
	velero.Status.Phase = "Running"
	err = r.Status().Update(ctx, velero)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *VeleroReconciler) desiredDeployment(velero *backupv1alpha1.Velero) *appsv1.Deployment {
	labels := map[string]string{
		"app": "velero",
	}
	replicas := int32(1)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "velero",
			Namespace: velero.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "velero",
						Image: "velero/velero:v1.7.1", // Use the appropriate Velero image version
						Env: []corev1.EnvVar{
							{
								Name:  "VELERO_NAMESPACE",
								Value: velero.Namespace,
							},
							{
								Name:  "VELERO_PROVIDER",
								Value: velero.Spec.Provider,
							},
							{
								Name:  "VELERO_BUCKET",
								Value: velero.Spec.Bucket,
							},
							{
								Name:  "VELERO_REGION",
								Value: velero.Spec.Region,
							},
						},
					}},
				},
			},
		},
	}
}

func (r *VeleroReconciler) deploymentEqual(d1, d2 *appsv1.Deployment) bool {
	// Compare the deployment specs to determine if an update is needed
	return d1.Spec.Replicas == d2.Spec.Replicas &&
		d1.Spec.Template.Spec.Containers[0].Image == d2.Spec.Template.Spec.Containers[0].Image &&
		compareEnvVars(d1.Spec.Template.Spec.Containers[0].Env, d2.Spec.Template.Spec.Containers[0].Env)
}

func compareEnvVars(env1, env2 []corev1.EnvVar) bool {
	if len(env1) != len(env2) {
		return false
	}
	envMap := make(map[string]string)
	for _, e := range env1 {
		envMap[e.Name] = e.Value
	}
	for _, e := range env2 {
		if envMap[e.Name] != e.Value {
			return false
		}
	}
	return true
}

func (r *VeleroReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&backupv1alpha1.Velero{}).
		Complete(r)
}
