/*
Copyright 2023.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "hl1.com/kubebuilder-example/api/v1"
)

// GuesstbookReconciler reconciles a Guesstbook object
type GuesstbookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.hl1.com,resources=guesstbooks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.hl1.com,resources=guesstbooks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.hl1.com,resources=guesstbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Guesstbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *GuesstbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	ctx = context.Background()
	// _ = r.Log.WithValues("guest book,", req.NamespacedName)
	obj := &webappv1.Guesstbook{}
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		fmt.Println("err is:", err.Error())
		log.Log.Info("ERR IS", err.Error())
	} else {
		fmt.Println("获取到数据:", obj.Spec.FirstName, obj.Spec.LastName)
	}
	// 初始化cr的status位running
	obj.Status.Status1 = "running"
	if err := r.Status().Update(ctx, obj, nil); err != nil {
		fmt.Println("无法更新状态")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GuesstbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guesstbook{}).
		Complete(r)
}
