# What is it

This is an operator-demo. You can use it to create CRD, and listen the events that contains
ADD/UPDATE/DELETE event. 

Event ADD: Create apps of nginx, k8s cluster can access it through cluster ip

Event UPDATE: Update the version of nginx 

Event DELETE: Delete all of pods, deployments, and services

# How to get it 

```
    
    cd $GOPATH/src/
    
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
 
 ## Run it on other namespace
 1. Create service account
```
    kubectl yaml/demo/service_account.yaml -n yourNamespace
```
 2. Create cluster role 
```
    kubectl yaml/demo/cluster_role.yaml
```
 3. Create cluster role binding
```
    kubectl yaml/demo/cluster_role_binding.yaml
```
 
Then you need update the file(yaml/demo/qiniu_nginx.yaml), add the serviceAccountName in it.
    
 ```
    serviceAccountName: yourServiceAccountName
 ```
Create the app
 ```
    kubectl create yaml/demo/qiniu_nginx.yaml -n yourNamespace
```
    


