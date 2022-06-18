#!/bin/bash

# 构建镜像
docker build -t lfdockerhub/httpserver:v1.0.0 .

# 推送到镜像仓库
docker push lfdockerhub/httpserver:v1.0.0

# 创建configMap
kubectl create -f configMap.yaml

# 创建deployment
kubectl crate -f httpserver.yaml
