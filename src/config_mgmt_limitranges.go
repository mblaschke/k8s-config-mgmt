package main

import (
	"k8s.io/api/core/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceLimitRanges(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("LimitRanges")

	existingList, err := mgmt.K8sService.ListLimitRanges(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.LimitRanges {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.LimitRange).DeepCopyInto(&k8sObject)

			mgmt.K8sService.UpdateLimitRange(namespace.Name, &k8sObject)
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)
			mgmt.K8sService.CreateLimitRange(namespace.Name, item.Object.(*v1.LimitRange))
		}
	}


	// cleanup
	if mgmt.Configuration.Config.LimitRanges.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.LimitRanges[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.DeleteLimitRange(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
