#!/usr/bin/env bash

image="$1"

echo $image

cd dockerfiles/qiniu-nginx/

GOOS=linux GOARCH=amd64 go install cmd/k8s/qiniu_nginx.go

cp $GOBIN/qiniu_nginx app

echo "build go app success"

docker build -t $image .

echo "build docker image: $image success"

docker push $image

rm -rf app
