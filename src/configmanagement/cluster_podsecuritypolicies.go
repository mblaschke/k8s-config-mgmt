package configmanagement

import (
	"k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementClusterPodSecurityPolicies struct {
	*K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) init() {
	mgmt.Logger.SubCategory("PodSecurityPolicies")
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) listExistingItems()  (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	objList, err := mgmt.K8sService.PodSecurityPolicies().List()

	for _, item := range objList {
		list[item.Name] = item.DeepCopyObject()
	}

	return list, err
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) listConfigItems() (map[string]config.ConfigObject) {
	return mgmt.clusterConfig.PodSecurityPolicies
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	configItem.(*v1beta1.PodSecurityPolicy).DeepCopyInto(k8sItem.(*v1beta1.PodSecurityPolicy))
	return &k8sItem
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) handleCreate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.PodSecurityPolicies().Create(k8sItem.(*v1beta1.PodSecurityPolicy))
	return err
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) handleUpdate(k8sItem runtime.Object) (error) {
	_, err := mgmt.K8sService.PodSecurityPolicies().Update(k8sItem.(*v1beta1.PodSecurityPolicy))
	return err
}

func (mgmt *K8sConfigManagementClusterPodSecurityPolicies) handleDelete(k8sItem runtime.Object) (error) {
	err := mgmt.K8sService.PodSecurityPolicies().Delete(k8sItem.(*v1beta1.PodSecurityPolicy).Name)
	return err
}
