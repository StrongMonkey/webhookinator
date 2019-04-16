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

	"github.com/rancher/rio/pkg/apis/webhookinator.rio.cattle.io/v1"
	clientset "github.com/rancher/rio/pkg/generated/clientset/versioned/typed/webhookinator.rio.cattle.io/v1"
	informers "github.com/rancher/rio/pkg/generated/informers/externalversions/webhookinator.rio.cattle.io/v1"
	listers "github.com/rancher/rio/pkg/generated/listers/webhookinator.rio.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type GitWebHookReceiverHandler func(string, *v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error)

type GitWebHookReceiverController interface {
	GitWebHookReceiverClient

	OnChange(ctx context.Context, name string, sync GitWebHookReceiverHandler)
	OnRemove(ctx context.Context, name string, sync GitWebHookReceiverHandler)
	Enqueue(namespace, name string)

	Cache() GitWebHookReceiverCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type GitWebHookReceiverClient interface {
	Create(*v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error)
	Update(*v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error)
	UpdateStatus(*v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.GitWebHookReceiver, error)
	List(namespace string, opts metav1.ListOptions) (*v1.GitWebHookReceiverList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.GitWebHookReceiver, err error)
}

type GitWebHookReceiverCache interface {
	Get(namespace, name string) (*v1.GitWebHookReceiver, error)
	List(namespace string, selector labels.Selector) ([]*v1.GitWebHookReceiver, error)

	AddIndexer(indexName string, indexer GitWebHookReceiverIndexer)
	GetByIndex(indexName, key string) ([]*v1.GitWebHookReceiver, error)
}

type GitWebHookReceiverIndexer func(obj *v1.GitWebHookReceiver) ([]string, error)

type gitWebHookReceiverController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.GitWebHookReceiversGetter
	informer          informers.GitWebHookReceiverInformer
	gvk               schema.GroupVersionKind
}

func NewGitWebHookReceiverController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.GitWebHookReceiversGetter, informer informers.GitWebHookReceiverInformer) GitWebHookReceiverController {
	return &gitWebHookReceiverController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromGitWebHookReceiverHandlerToHandler(sync GitWebHookReceiverHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.GitWebHookReceiver
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.GitWebHookReceiver))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *gitWebHookReceiverController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.GitWebHookReceiver))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateGitWebHookReceiverOnChange(updater generic.Updater, handler GitWebHookReceiverHandler) GitWebHookReceiverHandler {
	return func(key string, obj *v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, _ := updater(copyObj)
			if newObj != nil {
				copyObj = newObj.(*v1.GitWebHookReceiver)
			}
		}

		return copyObj, err
	}
}

func (c *gitWebHookReceiverController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *gitWebHookReceiverController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *gitWebHookReceiverController) OnChange(ctx context.Context, name string, sync GitWebHookReceiverHandler) {
	c.AddGenericHandler(ctx, name, FromGitWebHookReceiverHandlerToHandler(sync))
}

func (c *gitWebHookReceiverController) OnRemove(ctx context.Context, name string, sync GitWebHookReceiverHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromGitWebHookReceiverHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *gitWebHookReceiverController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *gitWebHookReceiverController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *gitWebHookReceiverController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *gitWebHookReceiverController) Cache() GitWebHookReceiverCache {
	return &gitWebHookReceiverCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *gitWebHookReceiverController) Create(obj *v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error) {
	return c.clientGetter.GitWebHookReceivers(obj.Namespace).Create(obj)
}

func (c *gitWebHookReceiverController) Update(obj *v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error) {
	return c.clientGetter.GitWebHookReceivers(obj.Namespace).Update(obj)
}

func (c *gitWebHookReceiverController) UpdateStatus(obj *v1.GitWebHookReceiver) (*v1.GitWebHookReceiver, error) {
	return c.clientGetter.GitWebHookReceivers(obj.Namespace).UpdateStatus(obj)
}

func (c *gitWebHookReceiverController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.GitWebHookReceivers(namespace).Delete(name, options)
}

func (c *gitWebHookReceiverController) Get(namespace, name string, options metav1.GetOptions) (*v1.GitWebHookReceiver, error) {
	return c.clientGetter.GitWebHookReceivers(namespace).Get(name, options)
}

func (c *gitWebHookReceiverController) List(namespace string, opts metav1.ListOptions) (*v1.GitWebHookReceiverList, error) {
	return c.clientGetter.GitWebHookReceivers(namespace).List(opts)
}

func (c *gitWebHookReceiverController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.GitWebHookReceivers(namespace).Watch(opts)
}

func (c *gitWebHookReceiverController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.GitWebHookReceiver, err error) {
	return c.clientGetter.GitWebHookReceivers(namespace).Patch(name, pt, data, subresources...)
}

type gitWebHookReceiverCache struct {
	lister  listers.GitWebHookReceiverLister
	indexer cache.Indexer
}

func (c *gitWebHookReceiverCache) Get(namespace, name string) (*v1.GitWebHookReceiver, error) {
	return c.lister.GitWebHookReceivers(namespace).Get(name)
}

func (c *gitWebHookReceiverCache) List(namespace string, selector labels.Selector) ([]*v1.GitWebHookReceiver, error) {
	return c.lister.GitWebHookReceivers(namespace).List(selector)
}

func (c *gitWebHookReceiverCache) AddIndexer(indexName string, indexer GitWebHookReceiverIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.GitWebHookReceiver))
		},
	}))
}

func (c *gitWebHookReceiverCache) GetByIndex(indexName, key string) (result []*v1.GitWebHookReceiver, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.GitWebHookReceiver))
	}
	return result, nil
}