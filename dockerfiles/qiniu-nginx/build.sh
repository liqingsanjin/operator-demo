#!/usr/bin/env bash

image="$1"

echo $image

cd dockerfiles/qiniu-nginx/

go build -o app $GOPATH/src/operator-demo/cmd/k8s/qiniu_nginx.go

echo "build go app success"

docker build -t $image .

echo "build docker image: $image success"

docker push $image
