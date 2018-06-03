package main

import (
	v13 "k8s.io/api/rbac/v1"
)

type K8sConfigManagementClusterClusterRoleBindings struct {
	K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterClusterRoleBindings) Manage() {
	mgmt.Logger.SubCategory("ClusterRoleBindings")

	cluster := mgmt.clusterConfig

	// check if anything is to do
	if !mgmt.Configuration.Config.ClusterRoleBindings.AutoCleanup && len(cluster.ClusterRoleBindings) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.ClusterRoleBindings().List()
	if err != nil {
		panic(err)
	}

	for _, item := range cluster.ClusterRoleBindings {
		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v13.ClusterRoleBinding).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ClusterRoleBindings().Update(&k8sObject)
				mgmt.handleOperationState(err)
			}

		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.ClusterRoleBindings().Create(item.Object.(*v13.ClusterRoleBinding))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.ClusterRoleBindings.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := cluster.ClusterRoleBindings[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.ClusterRoleBindings().Delete(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
