package configmanagement

import (
	"k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementClusterStorageClasses struct {
	*K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterStorageClasses) init() {
	mgmt.Logger.SubCategory("StorageClasses")
}

func (mgmt *K8sConfigManagementClusterStorageClasses) listExistingItems()  (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.StorageClasses().List()

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementClusterStorageClasses) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.clusterConfig.StorageClasses
}

func (mgmt *K8sConfigManagementClusterStorageClasses) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1.StorageClass).DeepCopyInto(k8sItem.(*v1.StorageClass))
	return &k8sItem
}

func (mgmt *K8sConfigManagementClusterStorageClasses) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.StorageClasses().Create(k8sItem.(*v1.StorageClass))
	return err
}

func (mgmt *K8sConfigManagementClusterStorageClasses) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.StorageClasses().Update(k8sItem.(*v1.StorageClass))
	return err
}

func (mgmt *K8sConfigManagementClusterStorageClasses) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.StorageClasses().Delete(k8sItem.(*v1.StorageClass).Name)
	return err
}
