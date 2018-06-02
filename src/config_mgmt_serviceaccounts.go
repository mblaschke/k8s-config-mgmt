package main

import (
	"k8s.io/api/core/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceServiceAccounts(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("ServiceAccount")

	existingList, err := mgmt.K8sService.ListServiceAccounts(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.ServiceAccounts {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.ServiceAccount).DeepCopyInto(&k8sObject)

			mgmt.K8sService.UpdateServiceAccount(namespace.Name, &k8sObject)
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)
			mgmt.K8sService.CreateServiceAccount(namespace.Name, item.Object.(*v1.ServiceAccount))
		}
	}


	// cleanup
	if mgmt.Configuration.Config.ServiceAccounts.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.ServiceAccounts[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.DeleteServiceAccount(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
