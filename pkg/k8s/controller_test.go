package k8s

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestController(t *testing.T) {
	as := assert.New(t)

	var (
		QiniuNginxKind    = os.Getenv("QiniuNginxKind")
		QiniuNginxPlural  = os.Getenv("QiniuNginxPlural")
		QiniuNginxGroup   = os.Getenv("QiniuNginxGroup")
		QiniuNginxVersion = os.Getenv("QiniuNginxVersion")
		QiniuNginxName    = os.Getenv("QiniuNginxName")
		namespace         = os.Getenv("namespace")
	)
	kubeconfigpath := os.Getenv("HOME") + "/.kube/config"
	operator, err := NewOperator(kubeconfigpath)
	as.NoError(err)

	ctl := NewQiniuNginxController(operator.GetKubeClient(), operator.GetApiClient(), namespace)

	crdConfig := &CRDConfig{
		Name:    QiniuNginxName,
		Kind:    QiniuNginxKind,
		Plural:  QiniuNginxPlural,
		Group:   QiniuNginxGroup,
		Version: QiniuNginxVersion,
		Scope:   v1beta1.NamespaceScoped,
		Obj:     &QiniuNginx{},
		ObjList: &QiniuNginxList{},
	}

	as.NoError(operator.WatchEvent(context.TODO(), &WatchConfig{
		WatchNamespace: "",
		ResyncPeriod:   0,
		Handlers:       ctl,
		SchemeBuilder: func(scheme *runtime.Scheme) error {
			scheme.AddKnownTypes(
				schema.GroupVersion{
					Group:   QiniuNginxGroup,
					Version: QiniuNginxVersion,
				},
				&QiniuNginx{},
				&QiniuNginxList{},
			)
			metav1.AddToGroupVersion(scheme, schema.GroupVersion{
				Group:   QiniuNginxGroup,
				Version: QiniuNginxVersion,
			})
			return nil
		},
	}, crdConfig))

	time.Sleep(60 * time.Second)
}
