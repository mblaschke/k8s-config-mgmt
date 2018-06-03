package main

import (
	"k8s.io/api/storage/v1"
)

type K8sConfigManagementClusterStorageClasses struct {
	K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterStorageClasses) Manage() {
	mgmt.Logger.SubCategory("StorageClasses")

	cluster := mgmt.clusterConfig

	// check if anything is to do
	if !mgmt.Configuration.Config.StorageClasses.AutoCleanup && len(cluster.StorageClasses) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.StorageClasses().List()
	if err != nil {
		panic(err)
	}

	for _, item := range cluster.StorageClasses {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.StorageClass).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.StorageClasses().Update(&k8sObject)
				mgmt.handleOperationState(err)
			}

		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.StorageClasses().Create(item.Object.(*v1.StorageClass))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.StorageClasses.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := cluster.StorageClasses[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.StorageClasses().Delete(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
