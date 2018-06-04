package main

type K8sConfigManagementScopeCluster struct {
	K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementScopeCluster) ClusterRoles() (*K8sConfigManagementClusterClusterRoles) {
	obj := &K8sConfigManagementClusterClusterRoles{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) ClusterRoleBinding() (*K8sConfigManagementClusterClusterRoleBindings) {
	obj := &K8sConfigManagementClusterClusterRoleBindings{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) PodSecurityPolicies() (*K8sConfigManagementClusterPodSecurityPolicies) {
	obj := &K8sConfigManagementClusterPodSecurityPolicies{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) StorageClasses() (*K8sConfigManagementClusterStorageClasses) {
	obj := &K8sConfigManagementClusterStorageClasses{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) Namespaces() (*K8sConfigManagementClusterNamespaces) {
	obj := &K8sConfigManagementClusterNamespaces{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.funcs = obj
	return obj
}
