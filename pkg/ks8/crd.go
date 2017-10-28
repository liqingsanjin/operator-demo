package ks8

import (
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type CRDConfig struct {
	Name    string // format: Plural.Group
	Kind    string
	Plural  string
	Group   string
	Version string
	Scope   v1beta1.ResourceScope

	Obj           runtime.Object
	ObjList       runtime.Object
	SchemeBuilder func(*runtime.Scheme) error
}