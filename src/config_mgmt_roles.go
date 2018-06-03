package main

import (
	v13 "k8s.io/api/rbac/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceRoles(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("Roles")

	// check if anything is to do
	if !mgmt.Configuration.Config.Roles.AutoCleanup && len(namespace.Roles) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingList, err := mgmt.K8sService.Roles().List(namespace.Name)
	if err != nil {
		panic(err)
	}

	for _, item := range namespace.Roles {

		if item.Object.(*v13.Role).Namespace == "" {
			item.Object.(*v13.Role).Namespace = namespace.Name
		}

		if k8sObject, ok := existingList[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v13.Role).DeepCopyInto(&k8sObject)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.Roles().Update(namespace.Name, &k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)

			if mgmt.IsNotDryRun() {
				_, err := mgmt.K8sService.Roles().Create(namespace.Name, item.Object.(*v13.Role))
				mgmt.handleOperationState(err)
			}
		}
	}


	// cleanup
	if mgmt.Configuration.Config.Roles.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.Roles[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.Roles().Delete(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
