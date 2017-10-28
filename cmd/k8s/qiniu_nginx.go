package main

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/liqingsanjin/operator-demo/pkg/k8s"
)

func main() {
	var (
		kubeconfigpath    = os.Getenv("INCLUSTER")
		QiniuNginxKind    = os.Getenv("QiniuNginxKind")
		QiniuNginxPlural  =	os.Getenv("QiniuNginxPlural")
		QiniuNginxGroup   = os.Getenv("QiniuNginxGroup")
		QiniuNginxVersion = os.Getenv("QiniuNginxVersion")
		QiniuNginxName    = os.Getenv("QiniuNginxName")
	)
	if kubeconfigpath == "yes" {
		kubeconfigpath = ""
	} else {
		kubeconfigpath = os.Getenv("HOME") + "/.kube/config"
	}
	operator, err := k8s.NewOperator(kubeconfigpath)
	if err != nil {
		log.Fatal(err)
	}

	crdConfig := &k8s.CRDConfig{
		Name:    QiniuNginxName,
		Kind:    QiniuNginxKind,
		Plural:  QiniuNginxPlural,
		Group:   QiniuNginxGroup,
		Version: QiniuNginxVersion,
		Scope:   v1beta1.NamespaceScoped,
		Obj:     &k8s.QiniuNginx{},
		ObjList: &k8s.QiniuNginxList{},
	}

	err = operator.CreateCRD(crdConfig)
	defer func() {
		policy := metav1.DeletePropagationBackground
		operator.DeleteCRD(QiniuNginxName, &metav1.DeleteOptions{
			PropagationPolicy: &policy,
		})
		log.Println("delete success")
	}()
	if err != nil {
		log.Println(err)
	}
	log.Println("create success")

	namespace, err := k8s.GetCurrentNS()

	if err != nil {
		log.Println(err)
	}
	log.Println(namespace)
	controller := k8s.NewQiniuNginxController(operator.GetKubeClient(), operator.GetApiClient(), namespace)
	operator.WatchEvent(context.TODO(), &k8s.WatchConfig{
		WatchNamespace: "",
		ResyncPeriod:   0,
		Handlers:       controller,
		SchemeBuilder: func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				schema.GroupVersion{
					Group:   QiniuNginxGroup,
					Version: QiniuNginxVersion,
				},
				&k8s.QiniuNginx{},
				&k8s.QiniuNginxList{},
			)
			metav1.AddToGroupVersion(scheme, schema.GroupVersion{
				Group:   QiniuNginxGroup,
				Version: QiniuNginxVersion,
			})
			return nil
		},
	}, crdConfig)
	ch := make(chan int)
	ch <- 1
}
