package main

import (
	v13 "k8s.io/api/rbac/v1"
)

func (mgmt *K8sConfigManagement) ManageNamespaceRoles(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("Roles")

	existingList, err := mgmt.K8sService.ListRoles(namespace.Name)
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

			mgmt.K8sService.UpdateRole(namespace.Name, &k8sObject)
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)
			mgmt.K8sService.CreateRole(namespace.Name, item.Object.(*v13.Role))
		}
	}


	// cleanup
	if mgmt.Configuration.Config.Roles.AutoCleanup {
		for _, k8sObject := range existingList {
			if _, ok := namespace.Roles[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.DeleteRole(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
