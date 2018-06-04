package main

import (
	"k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/core/v1"
)

type K8sConfigManagementNamespacePodDisruptionBudgets struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) init() {
	mgmt.Logger.SubCategory("PodDisruptionBudgets")
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.PodDisruptionBudgets().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) listConfigItems() (map[string]cfgObject) {
	return mgmt.Namespace.PodDisruptionBudgets
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ConfigMap).DeepCopyInto(k8sItem.(*v1.ConfigMap))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.PodDisruptionBudgets().Create(mgmt.Namespace.Name, k8sItem.(*v1beta1.PodDisruptionBudget))
	return err
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.PodDisruptionBudgets().Update(mgmt.Namespace.Name, k8sItem.(*v1beta1.PodDisruptionBudget))
	return err
}

func (mgmt *K8sConfigManagementNamespacePodDisruptionBudgets) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.PodDisruptionBudgets().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap).Name)
	return err
}
