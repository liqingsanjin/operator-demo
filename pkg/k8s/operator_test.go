package k8s

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	CRDKind    = "MyCRD"
	CRDPlural  = strings.ToLower(CRDKind) + "s"
	CRDGroup   = "example.com"
	CRDVersion = "v1"
	CRDName    = CRDPlural + "." + CRDGroup
)


func TestOperator(t *testing.T) {
	as := assert.New(t)
	kubeConfigPath := os.Getenv("HOME") + "/.kube/config"


	operator, err := NewOperator(kubeConfigPath)
	as.NoError(err)
	crdData := &CRDConfig{
		Name:    CRDName,
		Kind:    CRDKind,
		Plural:  CRDPlural,
		Group:   CRDGroup,
		Version: CRDVersion,
		Scope:   v1beta1.NamespaceScoped,

		Obj:           &MyCRD{},
		ObjList:       &MyCRDList{},
	}
	as.NoError(operator.CreateCRD(crdData))
	defer func() {
		policy := metav1.DeletePropagationBackground
		options := &metav1.DeleteOptions{
			PropagationPolicy: &policy,
		}
		as.NoError(operator.DeleteCRD(crdData.Name, options))
		log.Println("delete success")
	}()
	log.Println("create success")

	watchConfig := &WatchConfig{
		WatchNamespace: "",
		ResyncPeriod:   0,
		Handlers:       &EventHandler{},
		SchemeBuilder: func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				schema.GroupVersion{
					Group:   CRDGroup,
					Version: CRDVersion,
				},
				&MyCRD{},
				&MyCRDList{},
			)
			metav1.AddToGroupVersion(scheme, schema.GroupVersion{
				Group:   CRDGroup,
				Version: CRDVersion,
			})
			return nil
		},
	}
	operator.WatchEvent(context.TODO(), watchConfig, crdData)
	time.Sleep(60 * time.Second)
}

type EventHandler struct {
}

func (e *EventHandler) OnAdd(obj interface{}) {
	log.Println("add a obj: ")
	log.Println(obj)
}
func (e *EventHandler) OnUpdate(oldObj, newObj interface{}) {
	log.Println("update a obj: ")
	log.Println("old is: ")
	log.Println(oldObj)
	log.Println("new is:")
	log.Println(newObj)
}
func (e *EventHandler) OnDelete(obj interface{}) {
	log.Println("delete a obj: ")
	log.Println(obj)
}

