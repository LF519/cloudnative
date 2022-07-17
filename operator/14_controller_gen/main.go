package main

import (
	"context"
	"fmt"
	"log"

	v1 "github.com/cloudnative/operator/14_controller_gen/pkg/apis/baiding.tech/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// config
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Fatalln(err)
	}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = v1.Codecs.WithoutConversion()
	config.GroupVersion = &v1.GroupVersion

	client, err := rest.RESTClientFor(config)
	if err != nil {
		log.Fatalln(err)
	}

	foo := v1.Foo{}
	err = client.Get().Namespace("default").Resource("foos").Name("crd-test").Do(context.TODO()).Into(&foo)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(foo.Spec.Name)
	fmt.Println(foo.Spec.Replicas)

	// 测试生成的deepcopy方法
	newObj := foo.DeepCopy()
	newObj.Spec.Name = "test2"

	fmt.Println(foo.Spec.Name)
	fmt.Println(foo.Spec.Replicas)
	fmt.Println(newObj.Spec.Name)
	fmt.Println(newObj.Spec.Replicas)

}
