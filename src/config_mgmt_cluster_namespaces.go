package main

import (
	"k8s.io/api/core/v1"
)

type K8sConfigManagementClusterNamespaces struct {
	K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementClusterNamespaces) Manage() {
	mgmt.Logger.Main("Manage Namespaces")

	// check if anything is to do
	if !mgmt.Configuration.Config.Namespaces.AutoCleanup && len(mgmt.namespaces) == 0 {
		mgmt.Logger.Step("skipping")
		return
	}

	existingNamespaces, err := k8sService.Namespaces().List()
	if err != nil {
		panic(err)
	}

	// ensure namespace
	for _, item := range mgmt.namespaces {
		if k8sObject, ok := existingNamespaces[item.Name]; ok {
			mgmt.Logger.Step("Updating %v [labels:%v]", item.Name, item.Labels)

			k8sObject.Labels = item.Labels

			if mgmt.IsNotDryRun() {
				_, err := k8sService.Namespaces().Update(k8sObject)
				mgmt.handleOperationState(err)
			}
		} else {
			mgmt.Logger.Step("Create %v [labels:%v]", item.Name, item.Labels)

			k8sObject := &v1.Namespace{}
			k8sObject.Name = item.Name
			k8sObject.Namespace = item.Name
			k8sObject.Labels = item.Labels

			if mgmt.IsNotDryRun() {
				_, err := k8sService.Namespaces().Create(k8sObject)
				mgmt.handleOperationState(err)
			}
		}
	}

	// cleanup
	if mgmt.Configuration.Config.Namespaces.AutoCleanup {
		for _, k8sObject := range existingNamespaces {
			if _, ok := mgmt.namespaces[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.Namespaces().Delete(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
