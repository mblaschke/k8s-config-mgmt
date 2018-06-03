package main

import (
	"k8s.io/api/policy/v1beta1"
)

func (mgmt *K8sConfigManagement) ManageNamespacePodDisruptionBudgets(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("PodDisruptionBudgets")

	// check if anything is to do
	if !mgmt.Configuration.Config.PodDisruptionBudgets.AutoCleanup && len(namespace.PodDisruptionBudgets) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.PodDisruptionBudgets().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.PodDisruptionBudgets {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1beta1.PodDisruptionBudget).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.PodDisruptionBudgets().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.PodDisruptionBudgets().Create(namespace.Name, item.Object.(*v1beta1.PodDisruptionBudget))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.PodDisruptionBudgets.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.PodDisruptionBudgets[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.PodDisruptionBudgets().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
