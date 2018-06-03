package main

type K8sConfigManagementScopeCluster struct {
	K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementScopeCluster) ClusterRoles() (*K8sConfigManagementClusterClusterRoles) {
	return &K8sConfigManagementClusterClusterRoles{K8sConfigManagementBaseCluster: mgmt.K8sConfigManagementBaseCluster}
}

func (mgmt *K8sConfigManagementScopeCluster) ClusterRoleBinding() (*K8sConfigManagementClusterClusterRoleBindings) {
	return &K8sConfigManagementClusterClusterRoleBindings{K8sConfigManagementBaseCluster: mgmt.K8sConfigManagementBaseCluster}
}

func (mgmt *K8sConfigManagementScopeCluster) PodSecurityPolicies() (*K8sConfigManagementClusterPodSecurityPolicies) {
	return &K8sConfigManagementClusterPodSecurityPolicies{K8sConfigManagementBaseCluster: mgmt.K8sConfigManagementBaseCluster}
}

func (mgmt *K8sConfigManagementScopeCluster) StorageClasses() (*K8sConfigManagementClusterStorageClasses) {
	return &K8sConfigManagementClusterStorageClasses{K8sConfigManagementBaseCluster: mgmt.K8sConfigManagementBaseCluster}
}

func (mgmt *K8sConfigManagementScopeCluster) Namespaces() (*K8sConfigManagementClusterNamespaces) {
	return &K8sConfigManagementClusterNamespaces{K8sConfigManagementBaseCluster: mgmt.K8sConfigManagementBaseCluster}
}
