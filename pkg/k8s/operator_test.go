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
	testCRDKind    = "MyCRD"
	testCRDPlural  = strings.ToLower(testCRDKind) + "s"
	testCRDGroup   = "example.com"
	testCRDVersion = "v1"
	testCRDName    = testCRDPlural + "." + testCRDGroup
)

type MyCRD struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   MyCRDSpec   `json:"spec"`
	Status MyCRDStatus `json:"status"`
}

type MyCRDList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MyCRD `json:"items"`
}

type MyCRDSpec struct {
	Name string `json:"name"`
}

type MyCRDStatus struct {
	Msg string `json:"msg"`
}

func AddKnownTypesMyCRD(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		schema.GroupVersion{
			Group:   testCRDGroup,
			Version: testCRDVersion,
		},
		&MyCRD{},
		&MyCRDList{},
	)
	metav1.AddToGroupVersion(scheme, schema.GroupVersion{
		Group:   testCRDGroup,
		Version: testCRDVersion,
	})
	return nil
}

func TestNewOperator(t *testing.T) {
	as := assert.New(t)
	config := &OperatorConfig{
		KubeConfigPath: os.Getenv("HOME") + "/.kube/config",
		WatchNamespace: "",
		ResyncPeriod:   0,
		Handlers:       &EventHandler{},
		IsInCluster:    false,
	}

	operator, err := NewOperator(config)
	as.NoError(err)
	crdData := &CRDConfig{
		Name:    testCRDName,
		Kind:    testCRDKind,
		Plural:  testCRDPlural,
		Group:   testCRDGroup,
		Version: testCRDVersion,
		Scope:   v1beta1.NamespaceScoped,

		Obj:           &MyCRD{},
		ObjList:       &MyCRDList{},
		SchemeBuilder: AddKnownTypesMyCRD,
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

	operator.WatchEvents(context.TODO(), crdData)
	time.Sleep(100 * time.Second)
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
