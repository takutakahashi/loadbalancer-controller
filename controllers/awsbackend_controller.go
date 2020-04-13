/*


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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	loadbalancerv1beta1 "github.com/takutakahashi/loadbalancer-controller/api/v1beta1"
	"github.com/takutakahashi/loadbalancer-controller/pkg/terraform"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

// AWSBackendReconciler reconciles a AWSBackend object
type AWSBackendReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=loadbalancer.takutakahashi.dev,resources=awsbackends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=loadbalancer.takutakahashi.dev,resources=awsbackends/status,verbs=get;update;patch

func (r *AWSBackendReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("awsbackend", req.NamespacedName)
	var backend loadbalancerv1beta1.AWSBackend
	err := r.Get(ctx, req.NamespacedName, &backend)
	if apierrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}
	if backend.Status.Phase == "" {
		backend.Status.Phase = loadbalancerv1beta1.BackendPhaseProvisioning
	}
	if backend.ObjectMeta.Finalizers == nil || len(backend.ObjectMeta.Finalizers) == 0 {
		backend.ObjectMeta.Finalizers = []string{"loadbalancer.takutakahashi.dev"}
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if backend.ObjectMeta.DeletionTimestamp != nil {
		backend.Status.Phase = loadbalancerv1beta1.BackendPhaseDeleting
	}
	r.Update(ctx, &backend)
	return r.reconcile(ctx, backend)
}

func (r *AWSBackendReconciler) reconcile(ctx context.Context, backend loadbalancerv1beta1.AWSBackend) (ctrl.Result, error) {
	switch backend.Status.Phase {
	case loadbalancerv1beta1.BackendPhaseProvisioning:
		return r.ReconcileApply(ctx, backend)
	case loadbalancerv1beta1.BackendPhaseProvisioned:
		return r.ReconcileVerify(ctx, backend)
	case loadbalancerv1beta1.BackendPhaseDeleting:
		return r.ReconcileDelete(ctx, backend)
	case loadbalancerv1beta1.BackendPhaseDeleted:
		return ctrl.Result{}, r.Delete(ctx, &backend)
	default:
		return ctrl.Result{}, nil
	}
}
func (r *AWSBackendReconciler) ReconcileVerify(ctx context.Context, backend loadbalancerv1beta1.AWSBackend) (ctrl.Result, error) {
	r.Log.Info("Verify")
	tc, err := terraform.NewClientForAWSBackend(backend)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	status, err := tc.GetStatus()
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	backend.Status = status
	if backend.Status.Internal != true {
		reached, err := backend.ReachableAll()
		if err != nil || !reached {
			return ctrl.Result{Requeue: true}, err
		}
	}
	backend.Status.Phase = loadbalancerv1beta1.BackendPhaseReady
	return ctrl.Result{}, r.Update(ctx, &backend)
}

func (r *AWSBackendReconciler) ReconcileApply(ctx context.Context, backend loadbalancerv1beta1.AWSBackend) (ctrl.Result, error) {
	r.Log.Info("Apply")
	tc, err := terraform.NewClientForAWSBackend(backend)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = tc.Apply()
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}
	backend.Status.Phase = loadbalancerv1beta1.BackendPhaseProvisioned
	return ctrl.Result{}, r.Update(ctx, &backend)

}
func (r *AWSBackendReconciler) ReconcileDelete(ctx context.Context, backend loadbalancerv1beta1.AWSBackend) (ctrl.Result, error) {
	r.Log.Info("Delete")
	tc, err := terraform.NewClientForAWSBackend(backend)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = tc.Destroy()
	if err != nil {
		return ctrl.Result{}, err
	}
	backend.Finalizers = []string{}
	backend.Status.Phase = loadbalancerv1beta1.BackendPhaseDeleted
	return ctrl.Result{}, r.Update(ctx, &backend)
}

func (r *AWSBackendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&loadbalancerv1beta1.AWSBackend{}).
		Complete(r)
}
