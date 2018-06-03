package main

import (
	"k8s.io/api/core/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceConfigMaps(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("ConfigMaps")

	// check if anything is to do
	if !mgmt.Configuration.Config.ConfigMaps.AutoCleanup && len(namespace.ConfigMaps) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.ConfigMaps().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.ConfigMaps {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.ConfigMap).DeepCopyInto(&k8sObject)


			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ConfigMaps().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ConfigMaps().Create(namespace.Name, item.Object.(*v1.ConfigMap))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.ConfigMaps.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.ConfigMaps[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.ConfigMaps().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
