package main

import (
	"fmt"
	"informer-study/demo"
	"informer-study/k8sconfig"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

type PodHandler struct {
}

func (p PodHandler) OnAdd(obj interface{}, isInInitialList bool) {
	fmt.Println("OnAdd:", obj.(*v1.Pod).Name)
}

func (p PodHandler) OnUpdate(oldObj, newObj interface{}) {
	fmt.Println("OnUpdate:", newObj.(*v1.Pod).Name)
}

func (p PodHandler) OnDelete(obj interface{}) {
	fmt.Println("OnDelete:", obj.(*v1.Pod).Name)
}

var _ cache.ResourceEventHandler = &PodHandler{}

func main() {
	client := k8sconfig.InitClient()
	podListWatch := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	informer := demo.NewCustomInformer(podListWatch, &v1.Pod{}, &PodHandler{})
	informer.Run()
}
