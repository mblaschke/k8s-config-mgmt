package main

import (
		"k8s.io/api/settings/v1alpha1"
)

func (mgmt *K8sConfigManagement) ManageNamespacePodPresets(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("PodPresets")

	// check if anything is to do
	if !mgmt.Configuration.Config.PodPresets.AutoCleanup && len(namespace.PodPresets) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.PodPresets().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.PodPresets {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1alpha1.PodPreset).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.PodPresets().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.PodPresets().Create(namespace.Name, item.Object.(*v1alpha1.PodPreset))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.PodPresets.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.PodPresets[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.PodPresets().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
