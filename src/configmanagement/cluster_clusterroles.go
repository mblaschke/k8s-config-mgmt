package configmanagement

import (
	v13 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementClusterClusterRoles struct {
	*K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterClusterRoles) init() {
	mgmt.Logger.SubCategory("ClusterRoles")
}

func (mgmt *K8sConfigManagementClusterClusterRoles) listExistingItems()  (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.ClusterRoles().List()
	
	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}
	
	return list, err
}

func (mgmt *K8sConfigManagementClusterClusterRoles) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.clusterConfig.ClusterRoles
}

func (mgmt *K8sConfigManagementClusterClusterRoles) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v13.ClusterRole).DeepCopyInto(k8sItem.(*v13.ClusterRole))
	return &k8sItem
}

func (mgmt *K8sConfigManagementClusterClusterRoles) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ClusterRoles().Create(k8sItem.(*v13.ClusterRole))
	return err
}

func (mgmt *K8sConfigManagementClusterClusterRoles) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ClusterRoles().Update(k8sItem.(*v13.ClusterRole))
	return err
}

func (mgmt *K8sConfigManagementClusterClusterRoles) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.ClusterRoles().Delete(k8sItem.(*v13.ClusterRole).Name)
	return err
}
