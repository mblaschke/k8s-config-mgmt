package configmanagement

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/pkg/errors"
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagementClusterNamespaces struct {
	*K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterNamespaces) init() {
	mgmt.Logger.SubCategory("ClusterRoleBindings")
}

func (mgmt *K8sConfigManagementClusterNamespaces) listExistingItems()  (map[string]runtime.Object, error) {
	list := map[string]runtime.Object{}
	return list, errors.New("not implemented")
}

func (mgmt *K8sConfigManagementClusterNamespaces) listConfigItems() (map[string]config.ConfigObject) {
	return map[string]config.ConfigObject{}
}

func (mgmt *K8sConfigManagementClusterNamespaces) deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object) {
	return nil
}

func (mgmt *K8sConfigManagementClusterNamespaces) handleCreate(k8sItem runtime.Object) (error) {
	return errors.New("not implemented")
}

func (mgmt *K8sConfigManagementClusterNamespaces) handleUpdate(k8sItem runtime.Object) (error) {
	return errors.New("not implemented")
}

func (mgmt *K8sConfigManagementClusterNamespaces) handleDelete(k8sItem runtime.Object) (error) {
	return errors.New("not implemented")
}


func (mgmt *K8sConfigManagementClusterNamespaces) Manage() {
	mgmt.Logger.Category("Manage Namespaces")

	// check if anything is to do
	if !mgmt.Configuration.AutoCleanup && len(mgmt.namespaces) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingNamespaces, err := mgmt.K8sService.Namespaces().List()
	if err != nil {
		panic(err)
	}

	// ensure namespace
	for _, item := range mgmt.namespaces {
		if k8sObject, ok := existingNamespaces[item.Name]; ok {
			mgmt.Logger.Step("updating %v", item.Name)
			mgmt.Logger.StepResult("labels: %v", item.Labels)
			mgmt.Logger.StepResult("annotations: %v", item.Annotations)

			k8sObject.Labels = item.Labels
			k8sObject.Annotations = item.Annotations

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.Namespaces().Update(k8sObject)
				mgmt.handleOperationState(err)
			}
			statsNamespaces.updated++
		} else {
			mgmt.Logger.Step("create %v", item.Name)
			mgmt.Logger.StepResult("labels: %v", item.Labels)
			mgmt.Logger.StepResult("annotations: %v", item.Annotations)

			k8sObject := &v1.Namespace{}
			k8sObject.Name = item.Name
			k8sObject.Namespace = item.Name
			k8sObject.Labels = item.Labels
			k8sObject.Annotations = item.Annotations

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.Namespaces().Create(k8sObject)
				mgmt.handleOperationState(err)
			}
			statsNamespaces.created++
		}
	}

	// cleanup
	if mgmt.Configuration.AutoCleanup {
		for _, k8sObject := range existingNamespaces {
			if _, ok := mgmt.namespaces[k8sObject.Name]; !ok {
				mgmt.Logger.Step("deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := mgmt.K8sService.Namespaces().Delete(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
				statsNamespaces.deleted++
			}
		}
	}
}
