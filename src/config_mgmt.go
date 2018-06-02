package main

import (
		"k8s.io/api/core/v1"
	)

type K8sConfigManagement struct {
	Configuration Configuration
	K8sService Kubernetes
	Logger *DaemonLogger

	namespaces map[string]cfgNamespace
}


func (mgmt *K8sConfigManagement) Run() {
	mgmt.Init()
	mgmt.ManageNamespaces()
	mgmt.ManageConfiguration()
}

func (mgmt *K8sConfigManagement) Init() {
	var err error
	mgmt.namespaces, err = mgmt.Configuration.BuildNamespaceConfiguration()
	if err != nil {
		panic(err)
	}
}

func (mgmt *K8sConfigManagement) filter(name string, whitelist, blacklist []string) (bool) {
	return true
}

func (mgmt *K8sConfigManagement) IsNotDryRun() (run bool) {
	if !opts.DryRun {
		run = true
	} else {
		mgmt.Logger.StepResult("dry run")
	}
	return
}

func (mgmt *K8sConfigManagement) handleOperationState(err error) {
	if err == nil {
		mgmt.Logger.StepResult("ok")
	} else {
		mgmt.Logger.StepResult("failed [%v]", err)
	}
}

func (mgmt *K8sConfigManagement) ManageNamespaces() {
	mgmt.Logger.Main("Manage Namespaces")

	existingNamespaces, err := k8sService.ListNamespaces()
	if err != nil {
		panic(err)
	}

	// ensure namespace
	for _, item := range mgmt.namespaces {
		if k8sObject, ok := existingNamespaces[item.Name]; ok {
			mgmt.Logger.Step("Updating %v [labels:%v]", item.Name, item.Labels)

			k8sObject.Labels = item.Labels

			if mgmt.IsNotDryRun() {
				_, err := k8sService.UpdateNamespace(k8sObject)
				mgmt.handleOperationState(err)			}
		} else {
			mgmt.Logger.Step("Create %v [labels:%v]", item.Name, item.Labels)

			k8sObject := &v1.Namespace{}
			k8sObject.Name = item.Name
			k8sObject.Namespace = item.Name
			k8sObject.Labels = item.Labels

			if mgmt.IsNotDryRun() {
				_, err := k8sService.CreateNamespace(k8sObject)
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
					err := k8sService.DeleteNamespace(k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}

func (mgmt *K8sConfigManagement) ManageConfiguration() {
	mgmt.Logger.Main("Manage Configuration")

	for _, namespace := range mgmt.namespaces {
		mgmt.Logger.Category("Namespace %v", namespace.Name)

		mgmt.ManageConfigurationServiceAccount(namespace)
	}
}

func (mgmt *K8sConfigManagement) ManageConfigurationServiceAccount(namespace cfgNamespace) {
	mgmt.Logger.SubCategory("ServiceAccount")

	existingServiceAccounts, err := mgmt.K8sService.ListServiceAccounts(namespace.Name)
	if err != nil {
		panic(err)
	}
	
	for _, item := range namespace.ServiceAccounts {
		if k8sObject, ok := existingServiceAccounts[item.Name]; ok {
			mgmt.Logger.Step("Updating %v", item.Name)

			// update
			item.Object.(*v1.ServiceAccount).DeepCopyInto(&k8sObject)

			mgmt.K8sService.UpdateServiceAccount(namespace.Name, &k8sObject)
		} else {
			mgmt.Logger.Step("Creating %v", item.Name)
			mgmt.K8sService.CreateServiceAccount(namespace.Name, item.Object.(*v1.ServiceAccount))
		}
	}


	// cleanup
	if mgmt.Configuration.Config.ServiceAccounts.AutoCleanup {
		for _, k8sObject := range existingServiceAccounts {
			if _, ok := namespace.ServiceAccounts[k8sObject.Name]; !ok {
				mgmt.Logger.Step("Deleting %v", k8sObject.Name)

				if mgmt.IsNotDryRun() {
					err := k8sService.DeleteServiceAccount(namespace.Name, k8sObject.Name)
					mgmt.handleOperationState(err)
				}
			}
		}
	}
}
