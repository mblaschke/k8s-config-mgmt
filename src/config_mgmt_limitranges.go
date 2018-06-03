package main

import (
	"k8s.io/api/core/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceLimitRanges(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("LimitRanges")

	// check if anything is to do
	if !mgmt.Configuration.Config.LimitRanges.AutoCleanup && len(namespace.LimitRanges) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.LimitRanges().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.LimitRanges {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.LimitRange).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.LimitRanges().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.LimitRanges().Create(namespace.Name, item.Object.(*v1.LimitRange))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.LimitRanges.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.LimitRanges[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.LimitRanges().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
