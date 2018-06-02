package main

import (
		"k8s.io/api/core/v1"
)

type K8sConfigManagement struct {
	Configuration Configuration
	K8sService Kubernetes
	Logger *DaemonLogger
}


func (mgmt *K8sConfigManagement) Run() {
	mgmt.ManageNamespaces()
}

func (mgmt *K8sConfigManagement) ManageNamespaces() {
	mgmt.Logger.Main("Manage Namespaces")

	configuredNamespaces, err := mgmt.Configuration.Config.Namespaces.GetList()
	if err != nil {
		panic(err)
	}

	existingNamespaces, err := k8sService.ListNamespaces()
	if err != nil {
		panic(err)
	}

	// ensure namespace
	for _, configNamespace := range configuredNamespaces {
		if nsObject, ok := existingNamespaces[configNamespace.Name]; ok {
			mgmt.Logger.Step("Updating namespace %v [labels:%v]", configNamespace.Name, configNamespace.Labels)

			nsObject.Labels = configNamespace.Labels

			if !opts.DryRun {
				if _, err := k8sService.UpdateNamespace(nsObject); err == nil {
					mgmt.Logger.StepResult("ok")
				} else {
					mgmt.Logger.StepResult("failed [%v]", err)
				}
			} else {
				mgmt.Logger.StepResult("dry run")
			}
		} else {
			mgmt.Logger.Step("Create namespace %v [labels:%v]", configNamespace.Name, configNamespace.Labels)

			nsObject := v1.Namespace{}
			nsObject.Name = configNamespace.Name
			nsObject.Namespace = configNamespace.Name
			nsObject.Labels = configNamespace.Labels

			if !opts.DryRun {
				if _, err := k8sService.CreateNamespace(nsObject); err == nil {
					mgmt.Logger.StepResult("ok")
				} else {
					mgmt.Logger.StepResult("failed [%v]", err)
				}
			} else {
				mgmt.Logger.StepResult("dry run")
			}
		}
	}

	// cleanup
	if mgmt.Configuration.Config.Namespaces.AutoCleanup {
		for _, nsObject := range existingNamespaces {
			if _, ok := configuredNamespaces[nsObject.Name]; !ok {
				mgmt.Logger.Step("Delete namespace %v", nsObject.Name)

				if !opts.DryRun {
					if err := k8sService.DeleteNamespace(nsObject.Name); err == nil {
						mgmt.Logger.StepResult("ok")
					} else {
						mgmt.Logger.StepResult("failed [%v]", err)
					}
				} else {
					mgmt.Logger.StepResult("dry run")
				}
			}
		}
	}
}

func (mgmt *K8sConfigManagement) filter(name string, whitelist, blacklist []string) (bool) {
	return true
}
