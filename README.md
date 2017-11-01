# What is it

This is an operator-demo. You can use it to create CRD, and listen the events that contains
ADD/UPDATE/DELETE event. 

Case ADD : Create nginx's pods, k8s cluster can access it

Case UPDATE: Update the version of nginx 

Case DELETE: Delete all the pods


# How to use it

1. cd $operator-demo //enter the root dir of the project, $operator-demo is the project dir

2. sh dockerfiles/qiniu-nginx/build.sh $version //build app, docker image and push the image
        
3. Deploy the app in k8s cluster

```
kubectl create -f yaml/demo/qiniu_nginx.yaml
```

4. Create/Update/Delete an instance of the crd
```
kubectl create -f yaml/demo/test_qiniu_nginx.yaml

kubectl replace -f yaml/demo/test_qiniu_nginx.yaml

kubectl delete -f yaml/demo/test_qiniu_nginx.yaml
```


