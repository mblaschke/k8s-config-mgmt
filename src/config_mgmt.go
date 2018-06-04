package main

import (
	"os"
	)

type K8sConfigManagement struct {
	K8sConfigManagementBase
}

func (mgmt *K8sConfigManagement) Run() {
	mgmt.ManageConfiguration()
}

func (mgmt *K8sConfigManagement) Init() {
	var err error

	mgmt.clusterConfig, err = mgmt.Configuration.BuildClusterConfiguration()
	if err != nil {
		panic(err)
	}

	mgmt.namespaces, err = mgmt.Configuration.BuildNamespaceConfiguration()
	if err != nil {
		panic(err)
	}
}

func (mgmt *K8sConfigManagement) ManageConfiguration() {
	mgmt.Logger.Main("Manage Configuration")

	mgmt.Logger.Category("Cluster configuration")
	scope := mgmt.ClusterScope()
	scope.PodSecurityPolicies().Manage()
	scope.ClusterRoles().Manage()
	scope.ClusterRoleBinding().Manage()
	scope.StorageClasses().Manage()
	scope.Namespaces().Manage()
	os.Exit(1)

	for _, namespace := range mgmt.namespaces {
		mgmt.Logger.Category("Namespace %v", namespace.Name)

		scope := mgmt.NamespaceScope(namespace)
		scope.ServiceAccounts().Manage()
		scope.ConfigMaps().Manage()
		scope.Roles().Manage()
		scope.RoleBindings().Manage()
		scope.ResourceQuotas().Manage()
		scope.PodPresets().Manage()
		scope.PodDisruptionBudgets().Manage()
		scope.LimitRanges().Manage()
		scope.NetworkPolicies().Manage()
	}
}


