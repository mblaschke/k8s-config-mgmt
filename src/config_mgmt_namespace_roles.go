package main

import (
	v13 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/core/v1"
)

type K8sConfigManagementNamespaceRoles struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceRoles) init() {
	mgmt.Logger.SubCategory("Roles")
}

func (mgmt *K8sConfigManagementNamespaceRoles) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.Roles().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceRoles) listConfigItems() (map[string]cfgObject) {
	return mgmt.Namespace.Roles
}

func (mgmt *K8sConfigManagementNamespaceRoles) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ConfigMap).DeepCopyInto(k8sItem.(*v1.ConfigMap))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceRoles) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.Roles().Create(mgmt.Namespace.Name, k8sItem.(*v13.Role))
	return err
}

func (mgmt *K8sConfigManagementNamespaceRoles) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.Roles().Update(mgmt.Namespace.Name, k8sItem.(*v13.Role))
	return err
}

func (mgmt *K8sConfigManagementNamespaceRoles) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.Roles().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap).Name)
	return err
}
