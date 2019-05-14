/*
Copyright 2019 Rancher Labs.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v1

import (
	"context"

	"github.com/rancher/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	informers "k8s.io/client-go/informers/core/v1"
	clientset "k8s.io/client-go/kubernetes/typed/core/v1"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type PersistentVolumeClaimHandler func(string, *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)

type PersistentVolumeClaimController interface {
	PersistentVolumeClaimClient

	OnChange(ctx context.Context, name string, sync PersistentVolumeClaimHandler)
	OnRemove(ctx context.Context, name string, sync PersistentVolumeClaimHandler)
	Enqueue(namespace, name string)

	Cache() PersistentVolumeClaimCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type PersistentVolumeClaimClient interface {
	Create(*v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	Update(*v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	UpdateStatus(*v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.PersistentVolumeClaim, error)
	List(namespace string, opts metav1.ListOptions) (*v1.PersistentVolumeClaimList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.PersistentVolumeClaim, err error)
}

type PersistentVolumeClaimCache interface {
	Get(namespace, name string) (*v1.PersistentVolumeClaim, error)
	List(namespace string, selector labels.Selector) ([]*v1.PersistentVolumeClaim, error)

	AddIndexer(indexName string, indexer PersistentVolumeClaimIndexer)
	GetByIndex(indexName, key string) ([]*v1.PersistentVolumeClaim, error)
}

type PersistentVolumeClaimIndexer func(obj *v1.PersistentVolumeClaim) ([]string, error)

type persistentVolumeClaimController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.PersistentVolumeClaimsGetter
	informer          informers.PersistentVolumeClaimInformer
	gvk               schema.GroupVersionKind
}

func NewPersistentVolumeClaimController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.PersistentVolumeClaimsGetter, informer informers.PersistentVolumeClaimInformer) PersistentVolumeClaimController {
	return &persistentVolumeClaimController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromPersistentVolumeClaimHandlerToHandler(sync PersistentVolumeClaimHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.PersistentVolumeClaim
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.PersistentVolumeClaim))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *persistentVolumeClaimController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.PersistentVolumeClaim))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdatePersistentVolumeClaimOnChange(updater generic.Updater, handler PersistentVolumeClaimHandler) PersistentVolumeClaimHandler {
	return func(key string, obj *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.PersistentVolumeClaim)
			}
		}

		return copyObj, err
	}
}

func (c *persistentVolumeClaimController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *persistentVolumeClaimController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *persistentVolumeClaimController) OnChange(ctx context.Context, name string, sync PersistentVolumeClaimHandler) {
	c.AddGenericHandler(ctx, name, FromPersistentVolumeClaimHandlerToHandler(sync))
}

func (c *persistentVolumeClaimController) OnRemove(ctx context.Context, name string, sync PersistentVolumeClaimHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromPersistentVolumeClaimHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *persistentVolumeClaimController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *persistentVolumeClaimController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *persistentVolumeClaimController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *persistentVolumeClaimController) Cache() PersistentVolumeClaimCache {
	return &persistentVolumeClaimCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *persistentVolumeClaimController) Create(obj *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return c.clientGetter.PersistentVolumeClaims(obj.Namespace).Create(obj)
}

func (c *persistentVolumeClaimController) Update(obj *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return c.clientGetter.PersistentVolumeClaims(obj.Namespace).Update(obj)
}

func (c *persistentVolumeClaimController) UpdateStatus(obj *v1.PersistentVolumeClaim) (*v1.PersistentVolumeClaim, error) {
	return c.clientGetter.PersistentVolumeClaims(obj.Namespace).UpdateStatus(obj)
}

func (c *persistentVolumeClaimController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.PersistentVolumeClaims(namespace).Delete(name, options)
}

func (c *persistentVolumeClaimController) Get(namespace, name string, options metav1.GetOptions) (*v1.PersistentVolumeClaim, error) {
	return c.clientGetter.PersistentVolumeClaims(namespace).Get(name, options)
}

func (c *persistentVolumeClaimController) List(namespace string, opts metav1.ListOptions) (*v1.PersistentVolumeClaimList, error) {
	return c.clientGetter.PersistentVolumeClaims(namespace).List(opts)
}

func (c *persistentVolumeClaimController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.PersistentVolumeClaims(namespace).Watch(opts)
}

func (c *persistentVolumeClaimController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.PersistentVolumeClaim, err error) {
	return c.clientGetter.PersistentVolumeClaims(namespace).Patch(name, pt, data, subresources...)
}

type persistentVolumeClaimCache struct {
	lister  listers.PersistentVolumeClaimLister
	indexer cache.Indexer
}

func (c *persistentVolumeClaimCache) Get(namespace, name string) (*v1.PersistentVolumeClaim, error) {
	return c.lister.PersistentVolumeClaims(namespace).Get(name)
}

func (c *persistentVolumeClaimCache) List(namespace string, selector labels.Selector) ([]*v1.PersistentVolumeClaim, error) {
	return c.lister.PersistentVolumeClaims(namespace).List(selector)
}

func (c *persistentVolumeClaimCache) AddIndexer(indexName string, indexer PersistentVolumeClaimIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.PersistentVolumeClaim))
		},
	}))
}

func (c *persistentVolumeClaimCache) GetByIndex(indexName, key string) (result []*v1.PersistentVolumeClaim, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.PersistentVolumeClaim))
	}
	return result, nil
}
