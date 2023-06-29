package demo

import (
	"fmt"
	"informer-study/k8sconfig"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func List() {
	client := k8sconfig.InitClient()
	watchFromClient := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	list, err := watchFromClient.List(metaV1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T\n", list)
	podList := list.(*v1.PodList)
	for _, pod := range podList.Items {
		fmt.Println(pod.Name)
	}
}
