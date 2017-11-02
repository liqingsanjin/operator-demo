package k8s

import (
	log "github.com/sirupsen/logrus"
	apiserverclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensionsv1beta1 "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type CRDHandler struct {
}

func (c *CRDHandler) OnAdd(obj interface{}) {
	log.Printf("add: %v\n", obj)
}
func (c *CRDHandler) OnUpdate(oldObj, newObj interface{}) {
	log.Printf("update: %v\n", oldObj)
	log.Printf("updated: %v\n", newObj)
}
func (c *CRDHandler) OnDelete(obj interface{}) {
	log.Printf("delete: %v\n", obj)

}

type QiniuNginxController struct {
	kubeClient *kubernetes.Clientset
	apiClient  *apiserverclient.Clientset
	deploy     DeploymentInterface
	service    ServiceInterface
}

func NewQiniuNginxController(kubeClient *kubernetes.Clientset, apiClient *apiserverclient.Clientset, namespace string) *QiniuNginxController {
	controller := new(QiniuNginxController)
	controller.kubeClient = kubeClient
	controller.apiClient = apiClient
	controller.deploy = NewDeployments(kubeClient, namespace)
	controller.service = NewService(kubeClient, namespace)
	return controller
}

func (q *QiniuNginxController) OnAdd(obj interface{}) {
	qiniuNginx := obj.(*QiniuNginx)
	err := q.DeployNginx(qiniuNginx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("deploy success")
	}
}

func (q *QiniuNginxController) OnUpdate(oldObj, newObj interface{}) {
	qiniuNginx := newObj.(*QiniuNginx)
	err := q.UpdateNginxVersion(qiniuNginx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("update success")
	}
}

func (q *QiniuNginxController) OnDelete(obj interface{}) {
	qiniuNginx := obj.(*QiniuNginx)
	err := q.DeleteNginx(qiniuNginx)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("delete success")
	}
}

func (q *QiniuNginxController) DeployNginx(nginx *QiniuNginx) error {
	deploy := q.deploy
	_, err := deploy.Create(deploy.MakeConfig(&DeploymentData{
		Name: nginx.Spec.Name,
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: nginx.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx_deploy_test",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx_deploy_test",
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
	if err != nil {
		return err
	}
	_, err = q.service.Create(q.service.MakeConfig(&ServiceData{
		Name: nginx.Spec.Name,
		Labels: map[string]string{
			"app": "nginx_deploy_test",
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx_deploy_test",
			},
			Ports: []apiv1.ServicePort{
				apiv1.ServicePort{
					Port:     nginx.Spec.Port,
					Protocol: "TCP",
				},
			},
		},
	}))
	return err
}

func (q *QiniuNginxController) UpdateNginxVersion(newngxin *QiniuNginx) error {
	deploy := q.deploy
	_, err := deploy.Update(deploy.MakeConfig(&DeploymentData{
		Name: newngxin.Spec.Name,
		Spec: extensionsv1beta1.DeploymentSpec{
			Replicas: newngxin.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx_deploy_test",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx_deploy_test",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:  "nginx",
							Image: newngxin.Spec.Image,
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
	return err
}

func (q *QiniuNginxController) DeleteNginx(ngxin *QiniuNginx) error {
	policy := metav1.DeletePropagationBackground
	err := q.deploy.Delete(ngxin.Spec.Name, &metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		return err
	}

	return q.service.Delete(ngxin.Spec.Name)
}
