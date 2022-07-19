package utils

import (
	"bytes"
	"text/template"

	demoV1beta1 "github.com/kubebuilder-demo/api/v1beta1"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	networkingV1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func parseTemplate(templateName string, app *demoV1beta1.App) []byte {
	tmpl, err := template.ParseFiles("controllers/template/" + templateName + ".yaml")
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, app)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func NewDeployment(app *demoV1beta1.App) *appV1.Deployment {
	deployment := &appV1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", app), deployment)
	if err != nil {
		panic(err)
	}
	return deployment
}

func NewService(app *demoV1beta1.App) *coreV1.Service {
	service := &coreV1.Service{}
	err := yaml.Unmarshal(parseTemplate("service", app), service)
	if err != nil {
		panic(err)
	}
	return service
}

func NewIngress(app *demoV1beta1.App) *networkingV1.Ingress {
	ingress := &networkingV1.Ingress{}
	err := yaml.Unmarshal(parseTemplate("ingress", app), ingress)
	if err != nil {
		panic(err)
	}
	return ingress
}
