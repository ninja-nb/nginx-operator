/*
Copyright 2026.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	platformv1 "github.com/narendra/nginx-operator/api/v1"
)

// NginxClusterReconciler reconciles a NginxCluster object
type NginxClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	conditionAvailable   = "Available"
	conditionProgressing = "Progressing"
	conditionDegraded    = "Degraded"
)

// +kubebuilder:rbac:groups=platform.example.com,resources=nginxclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=platform.example.com,resources=nginxclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.example.com,resources=nginxclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NginxCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.24.1/pkg/reconcile
func (r *NginxClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	nginxCluster := &platformv1.NginxCluster{}
	if err := r.Get(ctx, req.NamespacedName, nginxCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxCluster.Name,
			Namespace: nginxCluster.Namespace,
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, deployment, func() error {
		labels := map[string]string{
			"app": nginxCluster.Name,
		}

		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: labels,
		}
		deployment.Spec.Template.Labels = labels
		deployment.Spec.Replicas = &nginxCluster.Spec.Replicas
		deployment.Spec.Template.Spec.Containers = []corev1.Container{
			{
				Name:  "nginx",
				Image: "nginx:1.27",
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: 80,
					},
				},
				ReadinessProbe: &corev1.Probe{
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path: "/",
							Port: intstr.FromInt(80),
						},
					},
					InitialDelaySeconds: 5,
					PeriodSeconds:       10,
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("100m"),
						corev1.ResourceMemory: resource.MustParse("128Mi"),
					},
					Limits: corev1.ResourceList{
						corev1.ResourceCPU:    resource.MustParse("500m"),
						corev1.ResourceMemory: resource.MustParse("256Mi"),
					},
				},
			},
		}

		return controllerutil.SetControllerReference(nginxCluster, deployment, r.Scheme)
	})
	if err != nil {
		log.Error(err, "Could not reconcile Deployment", "name", deployment.Name)
		return ctrl.Result{}, err
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxCluster.Name,
			Namespace: nginxCluster.Namespace,
		},
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, service, func() error {
		labels := map[string]string{
			"app": nginxCluster.Name,
		}

		service.Spec.Selector = labels
		service.Spec.Ports = []corev1.ServicePort{
			{
				Name:       "http",
				Port:       80,
				TargetPort: intstr.FromInt(80),
				Protocol:   corev1.ProtocolTCP,
			},
		}
		service.Spec.Type = corev1.ServiceTypeClusterIP

		return controllerutil.SetControllerReference(nginxCluster, service, r.Scheme)
	})
	if err != nil {
		log.Error(err, "Could not reconcile Service", "name", service.Name)
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, req.NamespacedName, deployment); err != nil {
		log.Error(err, "Could not fetch Deployment for status update", "name", deployment.Name)
		return ctrl.Result{}, err
	}

	desiredReplicas := nginxCluster.Spec.Replicas
	readyReplicas := deployment.Status.ReadyReplicas

	nginxCluster.Status.ReadyReplicas = readyReplicas
	available := readyReplicas == desiredReplicas && desiredReplicas > 0

	if available {
		apimeta.SetStatusCondition(&nginxCluster.Status.Conditions, metav1.Condition{
			Type:               conditionAvailable,
			Status:             metav1.ConditionTrue,
			Reason:             "DesiredReplicasReady",
			Message:            "Desired number of replicas is ready",
			ObservedGeneration: nginxCluster.Generation,
		})
		apimeta.SetStatusCondition(&nginxCluster.Status.Conditions, metav1.Condition{
			Type:               conditionProgressing,
			Status:             metav1.ConditionFalse,
			Reason:             "DesiredStateReached",
			Message:            "Reconciliation reached the desired state",
			ObservedGeneration: nginxCluster.Generation,
		})
	} else {
		apimeta.SetStatusCondition(&nginxCluster.Status.Conditions, metav1.Condition{
			Type:               conditionAvailable,
			Status:             metav1.ConditionFalse,
			Reason:             "ReplicasNotReady",
			Message:            "Waiting for pods to become ready",
			ObservedGeneration: nginxCluster.Generation,
		})
		apimeta.SetStatusCondition(&nginxCluster.Status.Conditions, metav1.Condition{
			Type:               conditionProgressing,
			Status:             metav1.ConditionTrue,
			Reason:             "Reconciling",
			Message:            "Resources are being reconciled",
			ObservedGeneration: nginxCluster.Generation,
		})
	}
	apimeta.SetStatusCondition(&nginxCluster.Status.Conditions, metav1.Condition{
		Type:               conditionDegraded,
		Status:             metav1.ConditionFalse,
		Reason:             "AsExpected",
		Message:            "No degraded condition detected",
		ObservedGeneration: nginxCluster.Generation,
	})

	if err := r.Status().Update(ctx, nginxCluster); err != nil {
		log.Error(err, "Could not update NginxCluster status", "name", nginxCluster.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1.NginxCluster{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named("nginxcluster").
		Complete(r)
}
