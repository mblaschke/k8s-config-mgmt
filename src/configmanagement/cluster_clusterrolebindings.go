package configmanagement

import (
	v13 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementClusterClusterRoleBindings struct {
	*K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) init() {
	mgmt.Logger.SubCategory("ClusterRoleBindings")
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) listExistingItems()  (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.ClusterRoleBindings().List()

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.clusterConfig.ClusterRoleBindings
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v13.ClusterRoleBinding).DeepCopyInto(k8sItem.(*v13.ClusterRoleBinding))
	return &k8sItem
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ClusterRoleBindings().Create(k8sItem.(*v13.ClusterRoleBinding))
	return err
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.ClusterRoleBindings().Update(k8sItem.(*v13.ClusterRoleBinding))
	return err
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.ClusterRoleBindings().Delete(k8sItem.(*v13.ClusterRoleBinding).Name)
	return err
}
