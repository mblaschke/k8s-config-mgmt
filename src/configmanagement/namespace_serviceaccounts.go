package configmanagement

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementNamespaceServiceAccounts struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) init() {
	mgmt.Logger.SubCategory("ServiceAccounts")
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.ServiceAccounts().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.Namespace.ServiceAccounts
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ServiceAccount).DeepCopyInto(k8sItem.(*v1.ServiceAccount))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ServiceAccounts().Create(mgmt.Namespace.Name, k8sItem.(*v1.ServiceAccount))
	return err
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ServiceAccounts().Update(mgmt.Namespace.Name, k8sItem.(*v1.ServiceAccount))
	return err
}

func (mgmt *K8sConfigManagementNamespaceServiceAccounts) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.ServiceAccounts().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ServiceAccount).Name)
	return err
}
