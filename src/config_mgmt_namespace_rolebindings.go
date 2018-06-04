package main

import (
	v13 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/api/core/v1"
)

type K8sConfigManagementNamespaceRoleBindings struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) init() {
	mgmt.Logger.SubCategory("RoleBindings")
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.RoleBindings().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) listConfigItems() (map[string]cfgObject) {
	return mgmt.Namespace.RoleBindings
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ConfigMap).DeepCopyInto(k8sItem.(*v1.ConfigMap))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.RoleBindings().Create(mgmt.Namespace.Name, k8sItem.(*v13.RoleBinding))
	return err
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.RoleBindings().Update(mgmt.Namespace.Name, k8sItem.(*v13.RoleBinding))
	return err
}

func (mgmt *K8sConfigManagementNamespaceRoleBindings) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.RoleBindings().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap).Name)
	return err
}
