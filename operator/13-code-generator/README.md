code-generator生成代码
```sh

# 代码生成的工作目录，也就是我们的项目路径
ROOT_PACKAGE="/root/go/src/github.com/cloudnative/operator/13-code-generator"
# API Group
CUSTOM_RESOURCE_NAME="crd.example.com"
# API Version
CUSTOM_RESOURCE_VERSION="v1"


cd /root/go/src/k8s.io/code-generator
# 执行代码自动生成，其中pkg/client是生成目标目录，pkg/apis是类型定义目录
./generate-groups.sh all "$ROOT_PACKAGE/pkg/client" "$ROOT_PACKAGE/pkg/apis" "$CUSTOM_RESOURCE_NAME:$CUSTOM_RESOURCE_VERSION"


/root/go/src/k8s.io/code-generator/generate-groups.sh all pkg/client pkg/apis crd.examle.com:v1 --go-header-file=/root/go/src/k8s.io/code-generator/hack/boilerplate.go.txt --output-base .
```