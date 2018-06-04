package configmanagement

import (
	v12 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementNamespaceNetworkPolicies struct {
	*K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) init() {
	mgmt.Logger.SubCategory("NetworkPolicies")
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) listExistingItems() (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.NetworkPolicies().List(mgmt.Namespace.Name)

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.Namespace.NetworkPolicies
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v12.NetworkPolicy).DeepCopyInto(k8sItem.(*v12.NetworkPolicy))
	return &k8sItem
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.NetworkPolicies().Create(mgmt.Namespace.Name, k8sItem.(*v12.NetworkPolicy))
	return err
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.NetworkPolicies().Update(mgmt.Namespace.Name, k8sItem.(*v12.NetworkPolicy))
	return err
}

func (mgmt *K8sConfigManagementNamespaceNetworkPolicies) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.NetworkPolicies().Delete(mgmt.Namespace.Name, k8sItem.(*v12.NetworkPolicy).Name)
	return err
}
