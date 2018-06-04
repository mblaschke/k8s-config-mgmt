package configmanagement

type K8sConfigManagementScopeNamespace struct {
	K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementScopeNamespace) ServiceAccounts() (*K8sConfigManagementNamespaceServiceAccounts) {
	obj := &K8sConfigManagementNamespaceServiceAccounts{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) ConfigMaps() (*K8sConfigManagementNamespaceConfigMaps) {
	obj := &K8sConfigManagementNamespaceConfigMaps{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) Roles() (*K8sConfigManagementNamespaceRoles) {
	obj := &K8sConfigManagementNamespaceRoles{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) RoleBindings() (*K8sConfigManagementNamespaceRoleBindings) {
	obj := &K8sConfigManagementNamespaceRoleBindings{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) ResourceQuotas() (*K8sConfigManagementNamespaceResourceQuotas) {
	obj := &K8sConfigManagementNamespaceResourceQuotas{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) PodPresets() (*K8sConfigManagementNamespacePodPresets) {
	obj := &K8sConfigManagementNamespacePodPresets{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) PodDisruptionBudgets() (*K8sConfigManagementNamespacePodDisruptionBudgets) {
	obj := &K8sConfigManagementNamespacePodDisruptionBudgets{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) LimitRanges() (*K8sConfigManagementNamespaceLimitRanges) {
	obj := &K8sConfigManagementNamespaceLimitRanges{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeNamespace) NetworkPolicies() (*K8sConfigManagementNamespaceNetworkPolicies) {
	obj := &K8sConfigManagementNamespaceNetworkPolicies{}
	obj.K8sConfigManagementBaseNamespace = &mgmt.K8sConfigManagementBaseNamespace
	obj.funcs = obj
	return obj
}

