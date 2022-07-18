1. kubebuilder init
```shell
kubebuilder init --domain baiding.tech
```

2. kuberbuilder create api
```shell
# 提示是否创建Resource和Controller时选y
kubebuilder create api --group ingress --version v1beta1 --kind App
```