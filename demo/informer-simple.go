package demo

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

//自己实现一个简单的informer

type CustomInformer struct {
	listWatch       *cache.ListWatch           //用于获取数据
	objType         runtime.Object             //用于反序列化
	resourceHandler cache.ResourceEventHandler //用于处理数据

	reflector *cache.Reflector //用于启动Reflector
	fIFO      *cache.DeltaFIFO //用于获取数据
	store     cache.Store      //用于存储数据
}

// 模拟informer
func NewCustomInformer(listWatch *cache.ListWatch, objType runtime.Object, resourceHandler cache.ResourceEventHandler) *CustomInformer {
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)
	fifo := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction:  cache.MetaNamespaceKeyFunc,
		KnownObjects: store,
	})
	reflector := cache.NewReflector(listWatch, objType, fifo, 0)

	return &CustomInformer{listWatch: listWatch, objType: objType, resourceHandler: resourceHandler, reflector: reflector, fIFO: fifo, store: store}
}

func (ci *CustomInformer) Run() {
	ch := make(chan struct{})
	go func() {
		ci.reflector.Run(ch)
	}()
	for {
		ci.fIFO.Pop(func(obj interface{}, isInInitialList bool) error {
			for _, delta := range obj.(cache.Deltas) {
				switch delta.Type {
				case cache.Sync, cache.Added:
					ci.store.Add(delta.Object)
					ci.resourceHandler.OnAdd(delta.Object, isInInitialList)
				case cache.Updated:
					if old, exists, err := ci.store.Get(delta.Object); err == nil && exists {
						ci.store.Update(delta.Object)
						ci.resourceHandler.OnUpdate(old, delta.Object)
					}
				case cache.Deleted:
					ci.store.Delete(delta.Object)
					ci.resourceHandler.OnDelete(delta.Object)
				}
			}
			return nil
		})
	}
}
