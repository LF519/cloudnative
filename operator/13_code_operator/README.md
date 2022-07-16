执行命令生成代码, 如果没有生成informers和listers, 加上参数`-v 10`可以查看报错信息
```shell
/root/go/src/k8s.io/code-generator/generate-groups.sh all github.com/cloudnative/operator/13_code_operator/pkg/client github.com/cloudnative/operator/13_code_operator/pkg/apis crd.example.com:v1 --go-header-file=/root/go/src/k8s.io/code-generator/hack/boilerplate.go.txt --output-base ../../../../
```