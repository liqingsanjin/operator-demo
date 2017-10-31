package k9s

import "github.com/emicklei/go-restful/log"

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
