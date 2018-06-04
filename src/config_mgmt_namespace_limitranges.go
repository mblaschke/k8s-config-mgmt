package main

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type K8sConfigManagementNamespaceLimitRanges struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) init() {
	mgmt.Logger.SubCategory("LimitRanges")
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.LimitRanges().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) listConfigItems() (map[string]cfgObject) {
	return mgmt.Namespace.LimitRanges
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ConfigMap).DeepCopyInto(k8sItem.(*v1.ConfigMap))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.LimitRanges().Create(mgmt.Namespace.Name, k8sItem.(*v1.LimitRange))
	return err
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.LimitRanges().Update(mgmt.Namespace.Name, k8sItem.(*v1.LimitRange))
	return err
}

func (mgmt *K8sConfigManagementNamespaceLimitRanges) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.LimitRanges().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap).Name)
	return err
}
