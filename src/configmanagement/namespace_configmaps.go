package configmanagement

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementNamespaceConfigMaps struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) init() {
	mgmt.Logger.SubCategory("ConfigMaps")
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.ConfigMaps().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.Namespace.ConfigMaps
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ConfigMap).DeepCopyInto(k8sItem.(*v1.ConfigMap))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ConfigMaps().Create(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap))
	return err
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ConfigMaps().Update(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap))
	return err
}

func (mgmt *K8sConfigManagementNamespaceConfigMaps) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.ConfigMaps().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ConfigMap).Name)
	return err
}
