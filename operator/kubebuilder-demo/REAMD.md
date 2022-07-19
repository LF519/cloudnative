### 开发
1. kubebuilder init
```shell
kubebuilder init --domain baiding.tech
```

2. kuberbuilder create api
```shell
# 提示是否创建Resource和Controller时选y
kubebuilder create api --group ingress --version v1beta1 --kind App
```

3. 创建相关的`webhook`
```shell
kubebuilder create webhook --group ingress --version v1beta1 --kind App --defaulting --programmatic-validation --conversion
```

4. 修改`typs.go`下的`AppSpec`, `app.controller.go`中添加操作kubenetes内建资源的权限,然后执行命令重新生成`rbac`, `crd`
```shell
make manifests
```

5. 编写controller的逻辑


6. 创建`crd`等资源
```shell
make install
```

7. 为了支持webhook, 需要打开 `config/default/kustomization.yaml`中的配置

8. 实现webhook的配置逻辑, 修改`api/v1beta1/app_webhook.go`文件

9. 测试`webhook`, 需要安装`cert-manager`来生成cert
```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.8.0/cert-manager.yaml
```


### 如何本地测试

1. 添加本地测试相关的代码

- config/dev

- Makefile

```shell
.PHONY: dev
dev: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	cd config/manager && $(KUSTOMIZE) edit set image controller=${IMG}
	$(KUSTOMIZE) build config/dev | kubectl apply -f -
.PHONY: undev
undev: manifests kustomize ## Deploy controller to the K8s cluster specified in ~/.kube/config.
	$(KUSTOMIZE) build config/dev | kubectl delete --ignore-not-found=$(ignore-not-found) -f -

```


2. 获取证书放到临时文件目录下

```shell
kubectl get secrets webhook-server-cert -n  kubebuilder-demo-system -o jsonpath='{..tls\.crt}' |base64 -d > certs/tls.crt
kubectl get secrets webhook-server-cert -n  kubebuilder-demo-system -o jsonpath='{..tls\.key}' |base64 -d > certs/tls.key
```

3. 修改main.go,让webhook server使用指定证书

```go
	if os.Getenv("ENVIRONMENT") == "DEV" {
		path, err := os.Getwd()
		if err != nil {
			setupLog.Error(err, "unable to get work dir")
			os.Exit(1)
		}
		options.CertDir = path + "/certs"
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
```
```shell
export ENVIRONMENT=DEV
```

4. 部署

```shell
make dev
```

5. 清理环境

```shell
make undev
```
