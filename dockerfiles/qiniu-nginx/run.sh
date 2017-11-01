#!/usr/bin/env bash

image="liqingsanjin/jingx:v$1"

echo $image

cd $GOPATH/src/operator-demo/dockerfiles/qiniu-nginx/

go build -o app $GOPATH/src/operator-demo/cmd/k8s/qiniu_nginx.go

echo "build go app success"

docker build -t $image .

echo "build docker image success"

docker push $image

kubectl create -f $GOPATH/src/operator-demo/yaml/demo/qiniu_nginx.yaml

echo "deploy image success"