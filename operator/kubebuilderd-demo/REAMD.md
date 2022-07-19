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

4. 修改`typs.go`下的`AppSpec`, 然后执行命令重新生成`rbac`, `crd`
```shell
make manifests
```

5. 编写controller的逻辑


6. 创建`crd`等资源
```shell
make install
```
