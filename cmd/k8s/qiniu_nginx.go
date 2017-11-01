package main

import (
	"context"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"operator-demo/pkg/k8s"
)

var (
	QiniuNginxKind    = "QiniuNginx"
	QiniuNginxPlural  = strings.ToLower(QiniuNginxKind) + "s"
	QiniuNginxGroup   = "example.com"
	QiniuNginxVersion = "v1"
	QiniuNginxName    = QiniuNginxPlural + "." + QiniuNginxGroup
	namespace         = "default"
)

func main() {
	kubeconfigpath := os.Getenv("INCLUSTER")
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
