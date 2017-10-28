package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type ServiceData struct {
	Name   string
	Labels map[string]string

	Spec apiv1.ServiceSpec
}

type ServiceInterface interface {
	MakeConfig(*ServiceData) *apiv1.Service
	Get(string) (*apiv1.Service, error)
	Create(*apiv1.Service) (*apiv1.Service, error)
	Delete(name string) error
}

type services struct {
	client corev1.ServiceInterface
}

func NewService(kclient *kubernetes.Clientset, ns string) ServiceInterface {
	return &services{
		client: kclient.CoreV1().Services(ns),
	}
}

func (s *services) Get(name string) (*apiv1.Service, error) {
	return s.client.Get(name, metav1.GetOptions{})
}

func (s *services) MakeConfig(rawData *ServiceData) *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   rawData.Name,
			Labels: rawData.Labels,
		},
		Spec: rawData.Spec,
	}
}

func (s *services) Create(config *apiv1.Service) (*apiv1.Service, error) {
	return s.client.Create(config)
}

func (s *services) Delete(name string) error {
	return s.client.Delete(name, nil)
}
