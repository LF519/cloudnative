#!/bin/bash

# 构建镜像
# docker build -t lfdockerhub/httpserver:v1.0.0-metrics .

# 推送到镜像仓库
# docker push lfdockerhub/httpserver:v1.0.0-metrics

# 创建configMap
kubectl create -f configMap.yaml

# 安装promethues, loki和grafana
helm repo add grafana https://grafana.github.io/helm-charts
helm upgrade --install loki grafana/loki-stack --set grafana.enabled=true,prometheus.enabled=true,prometheus.alertmanager.persistentVolume.enabled=false,prometheus.server.persistentVolume.enabled=false

# 部署deployment
kubectl crate -f httpserver.yaml
