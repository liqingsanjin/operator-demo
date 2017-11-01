#!/usr/bin/env bash

image="liqingsanjin/jingx:v$1"

cd ~/go/src/operator-demo/dockerfiles/qiniu-nginx/app

go build -o app ~/go/src/operator-demo/cmd/k8s/qiniu_nginx.go

docker build -t image .

docker push image

kubectl create -f ~/go/src/operator-demo/yaml/demo/qiniu_nginx.yaml