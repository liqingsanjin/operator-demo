#!/usr/bin/env bash

image="liqingsanjin/jingx:v$1"

cd $GOPATH/src/operator-demo/dockerfiles/qiniu-nginx/

go build -o app $GOPATH/src/operator-demo/cmd/k8s/qiniu_nginx.go

docker build -t image .

docker push image

kubectl create -f $GOPATH/src/operator-demo/yaml/demo/qiniu_nginx.yaml