package configmanagement

import "k8s-config-mgmt/src/config"

type K8sConfigManagementScopeNamespace struct {
	K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementScopeNamespace) ServiceAccounts(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceServiceAccounts) {
	obj := &K8sConfigManagementNamespaceServiceAccounts{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) ConfigMaps(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceConfigMaps) {
	obj := &K8sConfigManagementNamespaceConfigMaps{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) Roles(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceRoles) {
	obj := &K8sConfigManagementNamespaceRoles{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) RoleBindings(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceRoleBindings) {
	obj := &K8sConfigManagementNamespaceRoleBindings{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) ResourceQuotas(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceResourceQuotas) {
	obj := &K8sConfigManagementNamespaceResourceQuotas{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) PodPresets(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespacePodPresets) {
	obj := &K8sConfigManagementNamespacePodPresets{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) PodDisruptionBudgets(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespacePodDisruptionBudgets) {
	obj := &K8sConfigManagementNamespacePodDisruptionBudgets{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) LimitRanges(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceLimitRanges) {
	obj := &K8sConfigManagementNamespaceLimitRanges{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) NetworkPolicies(config config.ConfigurationManagementItem) (*K8sConfigManagementNamespaceNetworkPolicies) {
	obj := &K8sConfigManagementNamespaceNetworkPolicies{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.Configuration = config
	obj.funcs = obj
	return obj
}

