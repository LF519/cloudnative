/*
custom controller需要监听Service资源, 当Service发生变化时:
新增Service时:
	1. 包含指定的annotation, 创建Ingress资源对象
	2. 不包含指定annotation, 忽略
删除Service时:
	1. 删除Ingress资源对象

更新Service时:
	1. 包含指定的annotation, 检查Ingress资源对象是否存在, 不存在则新建; 存在则忽略
	2. 不包含指定的annotation, 检查Ingress资源对象是否存在, 存在则删除; 不存在则忽略
*/

package pkg

import (
	"context"
	"reflect"
	"time"

	coreApi "k8s.io/api/core/v1"
	netApi "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	informer "k8s.io/client-go/informers/core/v1"
	netInformer "k8s.io/client-go/informers/networking/v1"
	"k8s.io/client-go/kubernetes"
	coreLister "k8s.io/client-go/listers/core/v1"
	netLister "k8s.io/client-go/listers/networking/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	workNum  = 5
	maxRetry = 10
)

// controller数据结构
type controller struct {
	client        kubernetes.Interface
	ingressLister netLister.IngressLister
	serviceLister coreLister.ServiceLister
	queue         workqueue.RateLimitingInterface
}

// 定义添加service的eventHandler
func (c *controller) addService(obj interface{}) {
	c.enqueue(obj)
}

// 定义更新service的eventHandler
func (c *controller) updateService(oldObj interface{}, newObj interface{}) {
	// todo 比较annotation是不是相同
	if reflect.DeepEqual(oldObj, newObj) {
		return
	}
	c.enqueue(newObj)
}

// 定义删除ingress的eventHandler
func (c *controller) deleteIngress(obj interface{}) {
	ingress := obj.(*netApi.Ingress)
	OwnerReference := v1.GetControllerOf(ingress)

	if OwnerReference == nil {
		return
	}
	if OwnerReference.Kind != "Service" {
		return
	}

	c.queue.Add(ingress.Namespace + "/" + ingress.Name)

}

// 入workqueue方法
func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(err)
	}

	// 把对象的key放到workqueue里面
	c.queue.Add(key)
}

func (c *controller) worker() {
	for c.processNextItem() {
	}
}

func (c *controller) processNextItem() bool {
	item, shutdown := c.queue.Get()
	if shutdown {
		return false
	}

	defer c.queue.Done(item)

	key := item.(string)

	err := c.syncService(key)
	if err != nil {
		c.handlerError(key, err)
	}

	return true
}

// 错误处理
func (c *controller) handlerError(key string, err error) {
	if c.queue.NumRequeues(key) <= maxRetry {
		c.queue.AddRateLimited(key)
		return
	}
	runtime.HandleError(err)
	c.queue.Forget(key)
}

// 构造Ingress对象
func (c *controller) constructIngress(service *coreApi.Service) *netApi.Ingress {
	ingress := netApi.Ingress{}
	ingress.ObjectMeta.OwnerReferences = []v1.OwnerReference{
		*v1.NewControllerRef(service, coreApi.SchemeGroupVersion.WithKind("Service")),
	}
	ingress.Name = service.Name
	ingress.Namespace = service.Namespace

	pathType := netApi.PathTypePrefix
	icn := "nginx"
	ingress.Spec = netApi.IngressSpec{
		IngressClassName: &icn,
		Rules: []netApi.IngressRule{
			{
				Host: "example.com",
				IngressRuleValue: netApi.IngressRuleValue{
					HTTP: &netApi.HTTPIngressRuleValue{
						Paths: []netApi.HTTPIngressPath{
							{
								Path:     "/",
								PathType: &pathType,
								Backend: netApi.IngressBackend{
									Service: &netApi.IngressServiceBackend{
										Name: service.Name,
										Port: netApi.ServiceBackendPort{
											Number: 80,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return &ingress
}

func (c *controller) syncService(key string) error {
	namespaceKey, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// 删除
	service, err := c.serviceLister.Services(namespaceKey).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// 新增和删除
	_, ok := service.GetAnnotations()["ingress/http"]
	ingress, err := c.ingressLister.Ingresses(namespaceKey).Get(name)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	if ok && errors.IsNotFound(err) {
		// create ingress11
		ig := c.constructIngress(service)
		_, err := c.client.NetworkingV1().Ingresses(namespaceKey).Create(context.TODO(), ig, v1.CreateOptions{})
		if err != nil {
			return err
		}
	} else if !ok && ingress != nil {
		// delete ingress
		c.client.NetworkingV1().Ingresses(namespaceKey).Delete(context.TODO(), name, v1.DeleteOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// 启动controller
func (c *controller) Run(stopCh chan struct{}) {
	for i := 0; i < workNum; i++ {
		go wait.Until(c.worker, time.Minute, stopCh)
	}
	<-stopCh
}

// 构造一个controller实例
func NewController(client kubernetes.Interface, serviceInformer informer.ServiceInformer, ingeressInformer netInformer.IngressInformer) controller {
	c := controller{
		client:        client,
		ingressLister: ingeressInformer.Lister(),
		serviceLister: serviceInformer.Lister(),
		queue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ingressManager"),
	}

	serviceInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.addService,
		UpdateFunc: c.updateService,
	})

	ingeressInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: c.deleteIngress,
	})

	return c
}
