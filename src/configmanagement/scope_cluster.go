package configmanagement

type K8sConfigManagementScopeCluster struct {
	K8sConfigManagementBaseCluster
}

func (mgmt *K8sConfigManagementScopeCluster) ClusterRoles() (*K8sConfigManagementClusterClusterRoles) {
	obj := &K8sConfigManagementClusterClusterRoles{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.Configuration = mgmt.GlobalConfiguration.Management.Cluster.ClusterRoles
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) ClusterRolebindings() (*K8sConfigManagementClusterClusterRoleBindings) {
	obj := &K8sConfigManagementClusterClusterRoleBindings{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.Configuration = mgmt.GlobalConfiguration.Management.Cluster.ClusterRolebindings
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) PodSecurityPolicies() (*K8sConfigManagementClusterPodSecurityPolicies) {
	obj := &K8sConfigManagementClusterPodSecurityPolicies{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.Configuration = mgmt.GlobalConfiguration.Management.Cluster.PodSecurityPolicies
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) StorageClasses() (*K8sConfigManagementClusterStorageClasses) {
	obj := &K8sConfigManagementClusterStorageClasses{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.Configuration = mgmt.GlobalConfiguration.Management.Cluster.StorageClasses
	obj.funcs = obj
	return obj
}

func (mgmt *K8sConfigManagementScopeCluster) Namespaces() (*K8sConfigManagementClusterNamespaces) {
	obj := &K8sConfigManagementClusterNamespaces{}
	obj.K8sConfigManagementBaseCluster = &mgmt.K8sConfigManagementBaseCluster
	obj.Configuration = mgmt.GlobalConfiguration.Management.Cluster.Namespaces
	obj.funcs = obj
	return obj
}
