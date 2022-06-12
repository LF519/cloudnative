运行sh deploy.sh即可, 也可分别单独执行以下命令

* 构建镜像
    ```shell
    docker build -t lfdockerhub/httpserver:v1.0.0 .
    ```

* 推送到镜像仓库
    ```shell
    docker push lfdockerhub/httpserver:v1.0.0
    ```

* 创建configMap
    ```shell
    kubectl create -f configMap.yaml
    ```

* 创建deployment
    ```shell
    kubectl crate -f httpserver.yaml
    ```