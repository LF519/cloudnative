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

	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	demoV1Beta1 "github.com/kubebuilder-demo/api/v1beta1"
	"github.com/kubebuilder-demo/controllers/utils"
)

// AppReconciler reconciles a App object
type AppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=ingress.baiding.tech,resources=apps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=ingress.baiding.tech,resources=apps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=ingress.baiding.tech,resources=apps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the App object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *AppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// 从缓存中获取app
	app := &demoV1Beta1.App{}
	err := r.Get(ctx, req.NamespacedName, app)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 根据app的配置进行处理
	// 1. 处理Deployment
	// 	新建并绑定app与deployment的关系
	newDeployment := utils.NewDeployment(app)
	err = controllerutil.SetControllerReference(app, newDeployment, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	//	查找同名的deployment
	oldDeployment := &appV1.Deployment{}
	err = r.Get(ctx, req.NamespacedName, oldDeployment)
	if err != nil {
		// 若果找不到同名的deployment, 则新建一个
		if errors.IsNotFound(err) {
			err = r.Create(ctx, newDeployment)
			if err != nil {
				logger.Error(err, "create deploy failed")
				return ctrl.Result{}, err
			}
		}
	} else {
		// 如果同名的deployment已经存在, 更新
		err = r.Update(ctx, newDeployment)
		logger.Info("update deployment")
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// 2. Service的处理
	//	查找指定Service
	oldService := &coreV1.Service{}
	err = r.Get(ctx, req.NamespacedName, oldService)
	if err != nil {
		// 如果找不到指定的Service, 则新建一个
		if errors.IsNotFound(err) && app.Spec.EnableService {
			newService := utils.NewService(app)
			err = controllerutil.SetControllerReference(app, newService, r.Scheme)
			if err != nil {
				return ctrl.Result{}, err
			}
			err = r.Create(ctx, newService)
			if err != nil {
				logger.Error(err, "create service failed")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && app.Spec.EnableService {
			return ctrl.Result{}, err
		}
	} else {
		// 如果找到service, 并且EnableService的值为true, 跳过
		if app.Spec.EnableService {
			logger.Info("skip update service")
		} else {
			// 如果找到service, 并且EnableService的值为false, 需要删除
			err = r.Delete(ctx, oldService)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// 3. Ingress的处理, ingress的配置可能为空
	if !app.Spec.EnableService {
		return ctrl.Result{}, nil
	}

	oldIngress := &networkingV1.Ingress{}
	if err = r.Get(ctx, req.NamespacedName, oldIngress); err != nil {
		if errors.IsNotFound(err) && app.Spec.EnableIngress {
			newIngress := utils.NewIngress(app)
			if err = controllerutil.SetControllerReference(app, newIngress, r.Scheme); err != nil {
				return ctrl.Result{}, err
			}
			if err = r.Create(ctx, newIngress); err != nil {
				logger.Error(err, "create ingress failed")
				return ctrl.Result{}, err
			}
		}
		if !errors.IsNotFound(err) && app.Spec.EnableIngress {
			return ctrl.Result{}, err
		}
	} else {
		if app.Spec.EnableIngress {
			logger.Info("skip update ingress")
		} else {
			err = r.Delete(ctx, oldIngress)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *AppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demoV1Beta1.App{}).
		Complete(r)
}
