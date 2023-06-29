package demo

import (
	"fmt"
	"informer-study/k8sconfig"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func Watch() {
	client := k8sconfig.InitClient()
	watchFromClient := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	watch, err := watchFromClient.Watch(metaV1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for {
		select {
		case event, ok := <-watch.ResultChan():
			if !ok {
				return
			}
			fmt.Printf("%s:%s\n", event.Type, event.Object.(*v1.Pod).Name)
		}
	}
}
