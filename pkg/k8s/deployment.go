package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type DeploymentInterface interface {
	MakeConfig(*DeploymentData) *extensionsv1beta1.Deployment
	Create(*extensionsv1beta1.Deployment) (*extensionsv1beta1.Deployment, error)
	Get(string) (*extensionsv1beta1.Deployment, error)
	Delete(string, *metav1.DeleteOptions) error
	Update(*extensionsv1beta1.Deployment) (*extensionsv1beta1.Deployment, error)
}

type deployments struct {
	client    v1beta1.DeploymentInterface
	namespace string
}

type DeploymentData struct {
	Name   string
	Labels map[string]string

	Spec extensionsv1beta1.DeploymentSpec
}

func NewDeployments(kubeclient *kubernetes.Clientset, namespace string) DeploymentInterface {
	return &deployments{
		client:    kubeclient.ExtensionsV1beta1().Deployments(namespace),
		namespace: namespace,
	}
}

func (d *deployments) MakeConfig(data *DeploymentData) *extensionsv1beta1.Deployment {
	return &extensionsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: d.namespace,
			Labels:    data.Labels,
		},
		Spec: data.Spec,
	}
}

func (d *deployments) Create(deploy *extensionsv1beta1.Deployment) (*extensionsv1beta1.Deployment, error) {
	return d.client.Create(deploy)
}

func (d *deployments) Get(name string) (*extensionsv1beta1.Deployment, error) {
	return d.client.Get(name, metav1.GetOptions{})
}

func (d *deployments) Delete(name string, options *metav1.DeleteOptions) error {
	return d.client.Delete(name, options)
}

func (d *deployments) Update(deploy *extensionsv1beta1.Deployment) (*extensionsv1beta1.Deployment, error) {
	return d.client.Update(deploy)
}
