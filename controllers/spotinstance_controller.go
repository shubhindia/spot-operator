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

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// SpotInstanceReconciler reconciles a SpotInstance object
type SpotInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=shubhindia.xyz,resources=spotinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=shubhindia.xyz,resources=spotinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=shubhindia.xyz,resources=spotinstances/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpotInstance object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *SpotInstanceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpotInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Using v1.Pod here for faster testing. Since, both pod and node are part of core APIs
		// it will be easier to replace pod with node once I am done with initial logic
		For(&v1.Pod{}).
		Watches(
			&source.Kind{Type: &v1.Pod{}},
			handler.EnqueueRequestsFromMapFunc(r.markPreemptiveInstance),
		).
		WithEventFilter(
			predicate.Funcs{
				CreateFunc: func(ce event.CreateEvent) bool {
					return true
				},
				DeleteFunc: func(de event.DeleteEvent) bool {
					return true
				},
				UpdateFunc: func(ue event.UpdateEvent) bool {
					return false
				},
			},
		).
		Complete(r)
}

func (r *SpotInstanceReconciler) markPreemptiveInstance(o client.Object) []reconcile.Request {

	return nil
}
