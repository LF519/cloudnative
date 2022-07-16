* 自动生成yaml文件, 在生成的yaml文件上再做修改
    * deployment
        ```shell
        kubectl create deployment ingress-manager --image=lfdockerhub/ingress-manager:1.0.0 --namespace=default --dry-run=client -oyaml > manifests/ingress-manager.yaml
        ```

    * serviceaccount
        ```shell
        kubectl create sa ingress-manager-sa --namespace=default --dry-run=client -oyaml > manifests/ingress-manager-sa.yaml
        ```

    * role
        ```shell
        kubectl create clusterrole ingress-manager-role --resource=ingress,service --verb list,watch,create,update,delete --dry-run=client -oyaml > manifests/ingress-manager-role.yaml
        ```

    * rolebinding
        ```shell
        kubectl create clusterrolebinding ingress-manager-rb --clusterrole ingress-manager-role --serviceaccount default:ingress-manager-sa --dry-run=client -oyaml > manifests/ingress-manager-rb.yaml
        ```