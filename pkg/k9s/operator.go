package k9s

import (
	"time"

	"context"

	log "github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiserverclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type WatchConfig struct {
	WatchNamespace string
	ResyncPeriod   time.Duration
	Handlers       cache.ResourceEventHandler
	SchemeBuilder  func(*runtime.Scheme) error
}

type Operator struct {
	kubeClient     *kubernetes.Clientset
	apiClient      *apiserverclient.Clientset
	kubeconfig     *rest.Config
	operatorConfig *rest.Config
}

func NewOperator(kubeConfigPath string) (*Operator, error) {
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	o := new(Operator)
	o.apiClient, err = apiserverclient.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	o.kubeClient, err = kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	o.kubeconfig = kubeconfig
	return o, err
}

func (o *Operator) CreateCRD(crdconfig *CRDConfig) error {
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: crdconfig.Name,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   crdconfig.Group,
			Version: crdconfig.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: crdconfig.Plural,
				Kind:   crdconfig.Kind,
			},
		},
	}
	apiClientInterface := o.apiClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	crd, err := apiClientInterface.Create(crd)
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			_, err := apiClientInterface.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
			return err
		}
		return err
	}

	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err = apiClientInterface.Get(crd.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return true, err
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					log.Printf("Name conflict %v\n", cond.Reason)
				}

			}
		}
		return false, err
	})
	if err != nil {
		deleteErr := apiClientInterface.Delete(crd.Name, nil)
		if deleteErr != nil {
			return errors.NewAggregate([]error{err, deleteErr})
		}
		return err
	}
	return nil
}

func (o *Operator) WatchEvent(ctx context.Context, watchConfig *WatchConfig, crdConfig *CRDConfig) error {
	schemeBuilder := runtime.NewSchemeBuilder(watchConfig.SchemeBuilder)
	addToScheme := schemeBuilder.AddToScheme
	scheme := runtime.NewScheme()
	if err := addToScheme(scheme); err != nil {
		return err
	}

	config := o.kubeconfig
	config.GroupVersion = &schema.GroupVersion{
		Group:   crdConfig.Group,
		Version: crdConfig.Version,
	}
	config.APIPath = "/apis"
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
		crdConfig.Plural,
		watchConfig.WatchNamespace,
		fields.Everything(),
	)
	_, controller := cache.NewInformer(
		source,
		crdConfig.Obj,
		watchConfig.ResyncPeriod,
		watchConfig.Handlers,
	)
	go controller.Run(ctx.Done())
	return nil
}

func (o *Operator) DeleteCRD(name string, options *metav1.DeleteOptions) error {
	return o.apiClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(name, options)
}