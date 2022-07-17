下载并安装controller-tools
```shell
git clone git@github.com:kubernetes-sigs/controller-tools.git
cd controller-gen
go install ./cmd/{controller-gen,type-scaffold}
```

使用`type-scaffold`生成`typs.go`, 注意: 需要自己手工添加导包代码
```shell
~/go/bin/type-scaffold --kind Foo > pkg/apis/baiding.tech/v1/types.go
```

生成`deepcopy`方法文件`zz_generated.deepcopy.go`
```shell
~/go/bin/controller-gen object paths=pkg/apis/baiding.tech/v1/types.go
```

手工添加`register.go`文件, 代码如下
```go
// +groupName=baiding.tech
package v1
```

根据手工成的`register.go`文件, 执行命令生成`crd`
```shell
~/go/bin/controller-gen crd paths=./... output:crd:dir=config/crd
```