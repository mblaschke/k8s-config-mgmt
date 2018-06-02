package main

import (
	v13 "k8s.io/api/rbac/v1"
	)

func (mgmt *K8sConfigManagement) ManageClusterRoles() {
	mgmt.Logger.SubCategory("ClusterRoles")

	cluster := mgmt.clusterConfig

	// check if anything is to do
	if !mgmt.Configuration.Config.ClusterRoles.AutoCleanup && len(cluster.ClusterRoles) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.ListClusterRoles()
	if err != nil {
		panic(err)
	}

	for _, item := range cluster.ClusterRoles {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v13.ClusterRole).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.UpdateClusterRole(&k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.CreateClusterRole(item.Object.(*v13.ClusterRole))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.ClusterRoles.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := cluster.ClusterRoles[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.DeleteClusterRole(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
