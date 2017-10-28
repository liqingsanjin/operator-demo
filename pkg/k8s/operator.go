package k8s

import (
	"fmt"
	"time"

	"golang.org/x/net/context"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OperatorConfig struct {
	KubeConfigPath string

	WatchNamespace string

	ResyncPeriod time.Duration

	Handlers cache.ResourceEventHandler
}

type Operator struct {
	kubeconfig *rest.Config
	config     *OperatorConfig
	aclient    *clientset.Clientset
}

func NewOperator(config *OperatorConfig) (*Operator, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", config.KubeConfigPath)
	if err != nil {
		return nil, err
	}

	o := new(Operator)
	o.aclient, err = clientset.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	o.config = config
	o.kubeconfig = kubeconfig
	return o, nil
}

func (o *Operator) CreateCRD(crdconfig *CRDConfig) error {
	crdc := &v1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdconfig.Name,
		},
		Spec: v1beta1.CustomResourceDefinitionSpec{
			Group:   crdconfig.Group,
			Version: crdconfig.Version,
			Names: v1beta1.CustomResourceDefinitionNames{
				Kind:   crdconfig.Kind,
				Plural: crdconfig.Plural,
			},
			Scope: crdconfig.Scope,
		},
	}

	crdi := o.aclient.ApiextensionsV1beta1().CustomResourceDefinitions()
	crd, err := crdi.Create(crdc)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			_, err := crdi.Get(crdc.ObjectMeta.Name, metav1.GetOptions{})
			return err
		}
		return err
	}
	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err = crdi.Get(crdc.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case v1beta1.Established:
				if cond.Status == v1beta1.ConditionTrue {
					return true, err
				}
			case v1beta1.NamesAccepted:
				if cond.Status == v1beta1.ConditionFalse {
					fmt.Printf("Name conflict: %v\n", cond.Reason)
				}
			}
		}

		return false, err
	})
	if err != nil {
		if deleteErr := crdi.Delete(crdc.ObjectMeta.Name, nil); deleteErr != nil {
			return errors.NewAggregate([]error{err, deleteErr})
		}

		return err
	}
	return nil
}

func (o *Operator) WatchEvents(ctx context.Context, crd *CRDConfig) error {
	schemeBuilder := runtime.NewSchemeBuilder(crd.SchemeBuilder)
	addToScheme := schemeBuilder.AddToScheme

	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		return err
	}
	config := o.kubeconfig

	config.APIPath = "/apis"
	config.GroupVersion = &schema.GroupVersion{
		Group:   crd.Group,
		Version: crd.Version,
	}
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{
		CodecFactory: serializer.NewCodecFactory(scheme),
	}
	client, err := rest.RESTClientFor(config)
	if err != nil {
		return err
	}

	source := cache.NewListWatchFromClient(
		client,
		crd.Plural,
		o.config.WatchNamespace,
		fields.Everything(),
	)
	_, controller := cache.NewIndexerInformer(
		source,
		crd.Obj,
		o.config.ResyncPeriod,
		o.config.Handlers,
		cache.Indexers{},
	)

	go controller.Run(ctx.Done())

	return nil
}

func (o *Operator) DeleteCRD(name string, options *metav1.DeleteOptions) error {
	return o.aclient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(name, options)
}
