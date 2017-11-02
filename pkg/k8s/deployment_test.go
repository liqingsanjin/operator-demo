package k8s

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDeployments(t *testing.T) {
	as := assert.New(t)

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	as.NoError(err)
	kubeClient, err := kubernetes.NewForConfig(config)
	as.NoError(err)

	deploy := NewDeployments(kubeClient, "default")
	var count *int32
	*count = 1
	_, err = deploy.Create(deploy.MakeConfig(&DeploymentData{
		Name: "test_deploy",
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: count,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "test_deploy",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:  "nginx",
							Image: "nginx:latest",
							Ports: []apiv1.ContainerPort{
								apiv1.ContainerPort{
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}))
	as.NoError(err)
	defer func() {
		policy := metav1.DeletePropagationBackground
		err = deploy.Delete("test_deploy", &v1.DeleteOptions{
			PropagationPolicy: &policy,
		})
		as.NoError(err)
	}()
}
