package main

type K8sConfigManagementScopeNamespace struct {
	K8sConfigManagementBaseNamespace
}

func (mgmt *K8sConfigManagementScopeNamespace) ServiceAccounts() (*K8sConfigManagementNamespaceServiceAccounts) {
	return &K8sConfigManagementNamespaceServiceAccounts{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) ConfigMaps() (*K8sConfigManagementNamespaceConfigMaps) {
	return &K8sConfigManagementNamespaceConfigMaps{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) Roles() (*K8sConfigManagementNamespaceRoles) {
	return &K8sConfigManagementNamespaceRoles{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) RoleBindings() (*K8sConfigManagementNamespaceRoleBindings) {
	return &K8sConfigManagementNamespaceRoleBindings{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) ResourceQuotas() (*K8sConfigManagementNamespaceResourceQuotas) {
	return &K8sConfigManagementNamespaceResourceQuotas{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) PodPresets() (*K8sConfigManagementNamespacePodPresets) {
	return &K8sConfigManagementNamespacePodPresets{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) PodDisruptionBudgets() (*K8sConfigManagementNamespacePodDisruptionBudgets) {
	return &K8sConfigManagementNamespacePodDisruptionBudgets{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) LimitRanges() (*K8sConfigManagementNamespaceLimitRanges) {
	return &K8sConfigManagementNamespaceLimitRanges{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

func (mgmt *K8sConfigManagementScopeNamespace) NetworkPolicies() (*K8sConfigManagementNamespaceNetworkPolicies) {
	return &K8sConfigManagementNamespaceNetworkPolicies{K8sConfigManagementBaseNamespace: mgmt.K8sConfigManagementBaseNamespace}
}

