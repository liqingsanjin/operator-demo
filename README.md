# What is it

This is an operator-demo. You can use it to create CRD, and listen the events that contains
ADD/UPDATE/DELETE event. 

Case ADD : Create nginx's pods, k8s cluster can access it

Case UPDATE: Update the version of nginx 

Case DELETE: Delete all the pods

# How to get it 

```
    
    cd $GOPATH/src
    
    git clone https://github.com/liqingsanjin/operator-demo.git

```


# How to use it

1. Enter the root dir of the project, $operator-demo is the project dir

```

    cd $GOPATH/src/operator-demo
    

```

2. Build app, docker image and push the image
        
```

   sh dockerfiles/qiniu-nginx/build.sh yourImage

```

3. Update the the image in the file (yaml/demo/qiniu_nginx.yaml). Then deploy the app in k8s cluster

```

    kubectl create -f yaml/demo/qiniu_nginx.yaml

```

4. Create/Update/Delete an instance of the crd
```

    kubectl create -f yaml/demo/test_qiniu_nginx.yaml

    kubectl replace -f yaml/demo/test_qiniu_nginx.yaml

    kubectl delete -f yaml/demo/test_qiniu_nginx.yaml

```


