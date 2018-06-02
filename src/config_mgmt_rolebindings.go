package main

import (
	v13 "k8s.io/api/rbac/v1"
	)

func (mgmt *K8sConfigManagement) ManageNamespaceRoleBindings(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("RoleBindings")

	// check if anything is to do
	if !mgmt.Configuration.Config.RoleBindings.AutoCleanup && len(namespace.RoleBindings) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.ListRoleBindings(namespace.Name)
	if err != nil {
		panic(err)
	}

	for _, item := range namespace.RoleBindings {

		if item.Object.(*v13.RoleBinding).Namespace == "" {
			item.Object.(*v13.RoleBinding).Namespace = namespace.Name
		}

		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v13.RoleBinding).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.UpdateRoleBinding(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.CreateRoleBinding(namespace.Name, item.Object.(*v13.RoleBinding))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.RoleBindings.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.RoleBindings[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.DeleteRoleBinding(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
