package main

import (
	"k8s.io/api/settings/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/core/v1"
)

type K8sConfigManagementNamespacePodPresets struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespacePodPresets) init() {
	mgmt.Logger.SubCategory("PodPresets")
}

func (mgmt *K8sConfigManagementNamespacePodPresets) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.PodPresets().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespacePodPresets) listConfigItems() (map[string]cfgObject) {
	return mgmt.Namespace.PodPresets
}

func (mgmt *K8sConfigManagementNamespacePodPresets) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ConfigMap).DeepCopyInto(k8sItem.(*v1.ConfigMap))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespacePodPresets) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.PodPresets().Create(mgmt.Namespace.Name, k8sItem.(*v1alpha1.PodPreset))
	return err
}

func (mgmt *K8sConfigManagementNamespacePodPresets) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.PodPresets().Update(mgmt.Namespace.Name, k8sItem.(*v1alpha1.PodPreset))
	return err
}

func (mgmt *K8sConfigManagementNamespacePodPresets) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.PodPresets().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap).Name)
	return err
}

