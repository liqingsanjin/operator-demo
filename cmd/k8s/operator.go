package main

import (
	"context"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"operator-demo/pkg/k9s"
)

var (
	CRDKind    = "MyCRD"
	CRDPlural  = strings.ToLower(CRDKind) + "s"
	CRDGroup   = "example.com"
	CRDVersion = "v1"
	CRDName    = CRDPlural + "." + CRDGroup
)

func main() {
	kubeconfigpath := os.Getenv("KUBECONFIGPATH")
	//kubeconfigpath := os.Getenv("HOME") + "/.kube/config"
	operator, err := k9s.NewOperator(kubeconfigpath)
	if err != nil {
		log.Fatal(err)
	}

	crdConfig := &k9s.CRDConfig{
		Name:    CRDName,
		Kind:    CRDKind,
		Plural:  CRDPlural,
		Group:   CRDGroup,
		Version: CRDVersion,
		Scope:   v1beta1.NamespaceScoped,
		Obj:     &k9s.MyCRD{},
		ObjList: &k9s.MyCRDList{},
	}
	err = operator.CreateCRD(crdConfig)
	if err != nil {
		//log.Fatal(err)
	}
	defer func() {
		policy := metav1.DeletePropagationBackground
		operator.DeleteCRD(CRDName,  &metav1.DeleteOptions{
			PropagationPolicy: &policy,
		})
		log.Println("delete success")
	}()
	log.Println("create success")

	operator.WatchEvent(context.TODO(), &k9s.WatchConfig{
		WatchNamespace: "",
		ResyncPeriod:   0,
		Handlers:       &k9s.CRDHandler{},
		SchemeBuilder: func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				schema.GroupVersion{
					Group:   CRDGroup,
					Version: CRDVersion,
				},
				&k9s.MyCRD{},
				&k9s.MyCRDList{},
			)
			metav1.AddToGroupVersion(scheme, schema.GroupVersion{
				Group:   CRDGroup,
				Version: CRDVersion,
			})
			return nil
		},
	}, crdConfig)
	time.Sleep(10 * time.Second)
}
