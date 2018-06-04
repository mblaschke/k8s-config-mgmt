package configmanagement

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementNamespaceResourceQuotas struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) init() {
	mgmt.Logger.SubCategory("ResourceQuotas")
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.ResourceQuotas().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.Namespace.ResourceQuotas
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.ResourceQuota).DeepCopyInto(k8sItem.(*v1.ResourceQuota))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ResourceQuotas().Create(mgmt.Namespace.Name, k8sItem.(*v1.ResourceQuota))
	return err
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ResourceQuotas().Update(mgmt.Namespace.Name, k8sItem.(*v1.ResourceQuota))
	return err
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.ResourceQuotas().Delete(mgmt.Namespace.Name, k8sItem.(*v1.ResourceQuota).Name)
	return err
}
