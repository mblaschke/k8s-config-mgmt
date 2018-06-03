package main

import (
	"k8s.io/api/core/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceServiceAccounts(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("ServiceAccount")

	// check if anything is to do
	if !mgmt.Configuration.Config.ServiceAccounts.AutoCleanup && len(namespace.ServiceAccounts) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.ServiceAccounts().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.ServiceAccounts {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.ServiceAccount).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ServiceAccounts().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ServiceAccounts().Create(namespace.Name, item.Object.(*v1.ServiceAccount))
				mgmt.handleOperationState(err)
			}
		}
	}

	// cleanup
	if mgmt.Configuration.Config.ServiceAccounts.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.ServiceAccounts[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.ServiceAccounts().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
