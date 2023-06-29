package demo

import (
	"fmt"
	"informer-study/k8sconfig"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

//初识Indexer

func Indexer() {
	client := k8sconfig.InitClient()
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	watchFromClient := cache.NewListWatchFromClient(client.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction:  cache.MetaNamespaceKeyFunc,
		KnownObjects: store,
	})
	reflector := cache.NewReflector(watchFromClient, &v1.Pod{}, df, 0)
	ch := make(chan struct{})
	go func() {
		reflector.Run(ch)
	}()
	for {
		df.Pop(func(obj interface{}, isInInitialList bool) error {
			for _, delta := range obj.(cache.Deltas) {
				switch delta.Type {
				case cache.Sync, cache.Added:
					store.Add(delta.Object)
				case cache.Updated:
					store.Update(delta.Object)
				case cache.Deleted:
					store.Delete(delta.Object)
				}
				fmt.Println(delta.Type, ":", delta.Object.(*v1.Pod).Name)
			}
			return nil
		})
	}
}
