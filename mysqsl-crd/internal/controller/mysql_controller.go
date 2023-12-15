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
	"reflect"

	//corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "example.com/m/api/v1"
)

// MysqlReconciler reconciles a Mysql object
type MysqlReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=mysqls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=mysqls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.tutorial.kubebuilder.io,resources=mysqls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Mysql object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MysqlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	defer utilruntime.HandleCrash()
	logger := log.FromContext(ctx)
	logger.Info("receive reconcile event", "name", req.String())
	fmt.Println("req info:", req.Name, req.String(), req.NamespacedName)
	// 获取mysql对象
	logger.Info("get mysql object", "name", req.String())
	mysql := &batchv1.Mysql{}
	if err := r.Get(ctx, req.NamespacedName, mysql); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 如果mysql处于删除状态，则跳过
	if mysql.DeletionTimestamp != nil {
		logger.Info("mysql in deleting", "name", req.String())
		return ctrl.Result{}, nil
	}

	// 同步资源状态
	logger.Info("begin to sync mysql", "name", req.String())
	if err := r.syncMysql(ctx, mysql); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
func (r *MysqlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// 将 Reconcile 添加到 manager 中，这样当 manager 启动时它就会被启动
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 3,
		}).
		For(&batchv1.Mysql{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
	// 下面的是自带的
	//return ctrl.NewControllerManagedBy(mgr).
	//	For(&batchv1.Mysql{}).
	//	Complete(r)
}

const (
	mysqlLabelName = "tutorial.kubebuilder.io/mysql"
)

func (r *MysqlReconciler) syncMysql(ctx context.Context, obj *batchv1.Mysql) error {
	// 同步状态
	logger := log.FromContext(ctx)
	mysql := obj.DeepCopy()
	name := types.NamespacedName{
		Namespace: mysql.Namespace,
		Name:      mysql.Name,
	}
	owner := []metav1.OwnerReference{{
		APIVersion:         mysql.APIVersion,
		Kind:               mysql.Kind,
		Name:               mysql.Name,
		UID:                mysql.UID,
		Controller:         pointer.BoolPtr(true),
		BlockOwnerDeletion: pointer.BoolPtr(true),
	}}
	labels := map[string]string{
		mysqlLabelName: mysql.Name,
	}
	meta := metav1.ObjectMeta{
		Name:            mysql.Name,
		Namespace:       mysql.Namespace,
		Labels:          labels,
		OwnerReferences: owner,
	}
	deploy := &appsv1.Deployment{}
	if err := r.Get(ctx, name, deploy); err != nil {
		logger.Info("get deployment succeeded")
		if !errors.IsNotFound(err) {
			return err
		}
		deploy = &appsv1.Deployment{
			ObjectMeta: meta,
			Spec:       getDeploymentSpec(mysql, labels),
		}
		if err := r.Create(ctx, deploy); err != nil {
			return err
		}
		logger.Info("create deployment succeeded", "name", name.String())
	} else {
		want := getDeploymentSpec(mysql, labels)
		get := getSpecFromDeployment(deploy)
		if !reflect.DeepEqual(want, get) {
			newDeploy := deploy.DeepCopy()
			newDeploy.Spec = want
			if err := r.Update(ctx, newDeploy); err != nil {
				return err
			}
			logger.Info("update deployment success", "name", name.String())
		}
	}
	if *mysql.Spec.Replicas == deploy.Status.Replicas {
		mysql.Status.Code = batchv1.SuccessCode
	} else {
		mysql.Status.Code = batchv1.FailureCode
	}
	return r.Client.Status().Update(ctx, mysql)
}

func getDeploymentSpec(mysql *batchv1.Mysql, labels map[string]string) appsv1.DeploymentSpec {
	return appsv1.DeploymentSpec{
		Replicas: mysql.Spec.Replicas,
		Selector: metav1.SetAsLabelSelector(labels),
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{Labels: labels},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    "main",
						Image:   mysql.Spec.Image,
						Command: mysql.Spec.Command,
					},
				},
			},
		},
		Strategy:                appsv1.DeploymentStrategy{},
		MinReadySeconds:         0,
		RevisionHistoryLimit:    nil,
		Paused:                  false,
		ProgressDeadlineSeconds: nil,
	}
}

func getSpecFromDeployment(deploy *appsv1.Deployment) appsv1.DeploymentSpec {
	container := deploy.Spec.Template.Spec.Containers[0]
	return appsv1.DeploymentSpec{
		Replicas: deploy.Spec.Replicas,
		Selector: deploy.Spec.Selector,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: deploy.Spec.Template.Labels,
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:    container.Name,
						Image:   container.Image,
						Command: container.Command,
					},
				},
			},
		},
	}
}
