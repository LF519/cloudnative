
* 说明
    1. 随机延时代码已添加至`project/main.go`下的`rootHandler`方法里面
    2. Metric已经添加在`project/metrics/metrics.go`下
    3. 运行deploy.sh即可部署至k8s集群
    4. 从`Promethus`界面中查询延时指标数据, 见目录下的`promethus.png`
    5. 建一个`Grafana Dashboard`展现延时分配情况, 见目录下的`grafana.png`
