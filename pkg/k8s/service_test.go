package k8s

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func TestService(t *testing.T) {
	as := assert.New(t)

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	as.NoError(err)
	kubeClient, err := kubernetes.NewForConfig(config)
	as.NoError(err)

	deploy := NewService(kubeClient, "default")
	var count *int32
	*count = 1
	_, err = deploy.Create(deploy.MakeConfig(&ServiceData{
		Name: "test_service",
		Labels: map[string]string{
			"app": "test_service",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "test_service",
			},
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port:     80,
					Protocol: "TCP",
				},
			},
		},
	}))
	as.NoError(err)
	defer func() {
		err = deploy.Delete("test_service")
		as.NoError(err)
	}()
}
