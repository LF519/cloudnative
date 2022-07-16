package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// // config
	// config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	// if err != nil {
	// 	panic(err)
	// }
	// config.GroupVersion = &v1.SchemeGroupVersion
	// config.NegotiatedSerializer = scheme.Codecs
	// config.APIPath = "/api"

	// // client
	// restClient, err2 := rest.RESTClientFor(config)
	// if err2 != nil {
	// 	panic(err2)
	// }

	// // get data
	// pod := v1.Pod{}
	// err3 := restClient.Get().Namespace("default").Resource("pods").Name("test").Do(context.TODO()).Into(&pod)
	// if err3 != nil {
	// 	fmt.Println(err3)
	// } else {
	// 	fmt.Println(pod.Name)
	// }

	// configSet
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		panic(err)
	}

	clinetset, err2 := kubernetes.NewForConfig(config)
	if err2 != nil {
		panic(err2)
	}

	coreV1 := clinetset.CoreV1()
	pod, err3 := coreV1.Pods("istio-demo").Get(context.TODO(), "httpserver-b596cd979-j2j4q", v1.GetOptions{})
	if err3 != nil {
		fmt.Println(err3)
	} else {
		fmt.Println(pod)
	}
}
