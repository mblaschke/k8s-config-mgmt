package main

import (
	"k8s.io/api/policy/v1beta1"
)

func (mgmt *K8sConfigManagement) ManagePodSecurityPolicies() {
	mgmt.Logger.SubCategory("PodSecurityPolicies")

	cluster := mgmt.clusterConfig

	// check if anything is to do
	if !mgmt.Configuration.Config.PodSecurityPolicies.AutoCleanup && len(cluster.PodSecurityPolicies) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.PodSecurityPolicies().List()
	if err != nil {
		panic(err)
	}

	for _, item := range cluster.PodSecurityPolicies {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1beta1.PodSecurityPolicy).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.PodSecurityPolicies().Update(&k8sObject)
				mgmt.handleOperationState(err)
			}

		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.PodSecurityPolicies().Create(item.Object.(*v1beta1.PodSecurityPolicy))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.PodSecurityPolicies.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := cluster.PodSecurityPolicies[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.PodSecurityPolicies().Delete(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
