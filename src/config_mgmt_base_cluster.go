package main

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sConfigManagementBaseCluster struct {
	K8sConfigManagementBase
	funcs K8sConfigManagementBaseClusterFuncs
}

type K8sConfigManagementBaseClusterFuncs interface {
	init()
	listExistingItems() (map[string]runtime.Object, error)
	listConfigItems() (map[string]cfgObject)
	deepCloneObject(configItem, k8sItem runtime.Object) (*runtime.Object)
	handleCreate(k8sItem runtime.Object) (error)
	handleUpdate(k8sItem runtime.Object) (error)
	handleDelete(k8sItem runtime.Object) (error)
}

func (mgmt *K8sConfigManagementBaseCluster) Manage() {
	mgmt.funcs.init()

	existingList, err := mgmt.funcs.listExistingItems()
	if err != nil {
		panic(err)
	}

	configList := mgmt.funcs.listConfigItems()

	for _, item := range configList {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			updatedObject := mgmt.funcs.deepCloneObject(item.Object, k8sObject)

			if mgmt.IsNotDryRun() {
				mgmt.handleOperationState(mgmt.funcs.handleUpdate(*updatedObject))
			}

		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				mgmt.handleOperationState(mgmt.funcs.handleCreate(item.Object))
			}
		}
	}

	// cleanup
	if mgmt.Configuration.Config.ClusterRoleBindings.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := configList[k8sObject.(v1.Object).GetName()]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.(v1.Object).GetName())

				if mgmt.IsNotDryRun() {
					mgmt.handleOperationState(mgmt.funcs.handleDelete(k8sObject))
				}
			}
		}
	}
}
