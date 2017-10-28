package k8s

import (
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type CRDConfig struct {
	Kind    string
	Plural  string
	Name    string
	Group   string
	Version string
	Scope   apiextensionsv1beta1.ResourceScope
	Obj     runtime.Object
	ObjList runtime.Object
}

type QiniuNginx struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec QiniuNginxSpec `json:"spec"`
}

type QiniuNginxList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []QiniuNginx `json:"items"`
}

type QiniuNginxSpec struct {
	Replicas *int32 `json:"replicas"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Port     int32  `json:"port"`
}
