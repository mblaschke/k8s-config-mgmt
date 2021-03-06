package configmanagement

import (
	"k8s-config-mgmt/src/config"
)

type K8sConfigManagement struct {
	K8sConfigManagementBase
}

func (mgmt *K8sConfigManagement) Run() {
	mgmt.ManageConfiguration()
	mgmt.Statistics()
}

func (mgmt *K8sConfigManagement) Init() {
	var err error

	mgmt.clusterConfig, err = mgmt.GlobalConfiguration.BuildClusterConfiguration()
	if err != nil {
		panic(err)
	}

	mgmt.namespaces, err = mgmt.GlobalConfiguration.BuildNamespaceConfiguration()
	if err != nil {
		panic(err)
	}
}

func (mgmt *K8sConfigManagement) Statistics() {
	mgmt.Logger.Category("Statistics")
	mgmt.Logger.Step("namespaces")
	mgmt.Logger.StepResult("created:%v  updated:%v  recreated:%v  deleted:%v", statsNamespaces.created, statsNamespaces.updated, statsNamespaces.recreated, statsNamespaces.deleted)

	mgmt.Logger.Step("cluster objects")
	mgmt.Logger.StepResult("created:%v  updated:%v  recreated:%v  deleted:%v", statsClusterObjects.created, statsClusterObjects.updated, statsClusterObjects.recreated, statsClusterObjects.deleted)

	mgmt.Logger.Step("namespace objects")
	mgmt.Logger.StepResult("created:%v  updated:%v  recreated:%v  deleted:%v", statsNamespaceObjects.created, statsNamespaceObjects.updated, statsNamespaceObjects.recreated, statsNamespaceObjects.deleted)
}

func (mgmt *K8sConfigManagement) ManageConfiguration() {
	var serviceConfig config.ConfigurationManagementItem
	mgmt.Logger.Main("Manage Configuration")

	mgmt.Logger.Category("Cluster configuration")
	scope := mgmt.ClusterScope()

	serviceConfig = mgmt.GlobalConfiguration.Management.Cluster.PodSecurityPolicies
	if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
		scope.PodSecurityPolicies().Manage()
	}

	serviceConfig = mgmt.GlobalConfiguration.Management.Cluster.ClusterRoles
	if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
		scope.ClusterRoles().Manage()
	}

	serviceConfig = mgmt.GlobalConfiguration.Management.Cluster.ClusterRolebindings
	if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
		scope.ClusterRolebindings().Manage()
	}

	serviceConfig = mgmt.GlobalConfiguration.Management.Cluster.StorageClasses
	if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
		scope.StorageClasses().Manage()
	}

	serviceConfig = mgmt.GlobalConfiguration.Management.Cluster.Namespaces
	if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
		scope.Namespaces().Manage()
	}

	for _, namespace := range mgmt.namespaces {
		mgmt.Logger.Category("Namespace %v", namespace.Name)

		scope := mgmt.NamespaceScope(namespace)

		namespaceConfiguration := config.ConfigurationManagementNamespace{}
		for _, nsConfig := range mgmt.GlobalConfiguration.Management.Namespaces {
			nameRegexp := filterValueToRegexp(nsConfig.Name)
			if nameRegexp.MatchString(namespace.Name) {
				namespaceConfiguration = nsConfig
				break
			}
		}

		serviceConfig := namespaceConfiguration.ServiceAccounts
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.ServiceAccounts(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.ConfigMaps
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.ConfigMaps(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.Roles
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.Roles(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.RoleBindings
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.RoleBindings(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.ResourceQuotas
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.ResourceQuotas(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.PodPresets
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.PodPresets(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.PodDisruptionBudgets
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.PodDisruptionBudgets(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.LimitRanges
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.LimitRanges(serviceConfig).Manage()
		}

		serviceConfig = namespaceConfiguration.NetworkPolicies
		if serviceConfig.Enabled == nil || (*serviceConfig.Enabled) == true{
			scope.NetworkPolicies(serviceConfig).Manage()
		}
	}
}


