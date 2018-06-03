package main

import (
		v12 "k8s.io/api/networking/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceNetworkPolicies(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("NetworkPolicies")

	// check if anything is to do
	if !mgmt.Configuration.Config.NetworkPolicies.AutoCleanup && len(namespace.NetworkPolicies) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.NetworkPolicies().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.NetworkPolicies {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v12.NetworkPolicy).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.NetworkPolicies().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.NetworkPolicies().Create(namespace.Name, item.Object.(*v12.NetworkPolicy))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.NetworkPolicies.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.NetworkPolicies[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.NetworkPolicies().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
