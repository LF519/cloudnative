* 创建名为`istio-demo`的namespace, 并打上`istio-injection=enabled`标签
    ```sh
    kubectl create ns istio-demo
    kubectl label ns istio-demo istio-injection=enabled
    ```

* 创建服务
    ```
    kubectl create -f httpserver.yaml
    ```

* 生成证书
    ```sh
    openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=cncamp Inc./CN=*.cncamp.io' -keyout cncamp.io.key -out cncamp.io.crt

    kubectl create -n istio-system secret tls cncamp-credential --key=cncamp.io.key --cert=cncamp.io.crt
    ```
* 创建`Gateway`
    ```sh
    kubectl create -f httpgw.yaml
    ```

* 创建`virtualService`
    ```sh
    kubectl create -f vsserver.yaml
    ```