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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gcpUtils "github.com/shubhindia/spot-operator/controllers/utils/gcp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// SpotInstanceReconciler reconciles a SpotInstance object
type SpotInstanceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	log    logr.Logger
}

//+kubebuilder:rbac:groups=shubhindia.xyz,resources=spotinstances,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=shubhindia.xyz,resources=spotinstances/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=shubhindia.xyz,resources=spotinstances/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

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
	r.log = log.FromContext(ctx).WithValues("Spot-operator", "all")

	// get all nodes
	nodes := &v1.NodeList{}

	_ = r.Client.List(ctx, nodes)

	nodesForDeletetion := []v1.Node{}

	for _, node := range nodes.Items {

		// check for preemtible node label
		if node.Labels["cloud.google.com/gke-preemptible"] == "true" {

			// only cordon node if it was created 23 hours ago and is not already cordoned
			if time.Since(node.CreationTimestamp.Time) > 23*time.Minute && !node.Spec.Unschedulable {
				err := r.Client.Patch(ctx, &v1.Node{
					ObjectMeta: metav1.ObjectMeta{
						Name: node.Name,
					},
					Spec: v1.NodeSpec{
						Unschedulable: true,
					},
				}, client.Merge)

				if err != nil {
					r.log.Info(fmt.Sprintf("Unable to cordon node %s", node.Name))

				}
				r.log.Info(fmt.Sprintf("Successfully cordoned node %s", node.Name))

				// mark node for deletion since its already passed its 23 hours mark and is already cordoned
				nodesForDeletetion = append(nodesForDeletetion, node)
			}

		}
	}

	for _, node := range nodesForDeletetion {

		podList := &v1.PodList{}
		err := r.Client.List(ctx, podList, &client.ListOptions{
			Raw: &metav1.ListOptions{
				FieldSelector: "spec.nodeName=" + node.Name,
			},
		})
		if err != nil {
			r.log.Error(err, "Failed to get podList")
		}
		for _, pod := range podList.Items {

			if pod.Labels["podName"] == "high-config" || pod.Labels["podName"] == "mid-config" || pod.Labels["podName"] == "low-config" {
				r.log.Info(fmt.Sprintf("Runner pod exists on node: %s", node.Name))

			} else {

				// Draining a node hasn't been in implemented in go-client yet. So for now, we directly delete the node
				// as the code itself is tailored for our specific use-case.
				// TODO: Add custom function for draining the node first.
				r.log.Info(fmt.Sprintf("Deleting unerlying VM for node %s", node.Name))

				// Our end goal here is to reset the preemptible node clock
				// https://cloud.google.com/compute/docs/instances/preemptible#preemption-selection
				// Since, these nodes are part of a node-pool, we can directly delete them and cluster-autoscaler will take care of bringing the cluster up to desired state.

				// TODO: Maintain a cache about nodes which are deleted previously.
				err = gcpUtils.DeleteNode("mobile-ci-infra", "asia-east1-a", node.Name)
				if err != nil {
					r.log.Error(err, fmt.Sprintf("Failed to delete instance %s", node.Name))
				}
			}
		}

	}

	return ctrl.Result{
		RequeueAfter: 5 * time.Minute,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpotInstanceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Node{}).
		Complete(r)
}
