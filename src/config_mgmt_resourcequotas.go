package main

import (
	"k8s.io/api/core/v1"
)

type K8sConfigManagementNamespaceResourceQuotas struct {
	K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementNamespaceResourceQuotas) Manage() {
	mgmt.Logger.SubCategory("ResourceQuotas")

	namespace := mgmt.Namespace

	// check if anything is to do
	if !mgmt.Configuration.Config.ResourceQuotas.AutoCleanup && len(namespace.ResourceQuotas) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.ResourceQuotas().List(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.ResourceQuotas {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.ResourceQuota).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ResourceQuotas().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ResourceQuotas().Create(namespace.Name, item.Object.(*v1.ResourceQuota))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.ResourceQuotas.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.ResourceQuotas[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.ResourceQuotas().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
