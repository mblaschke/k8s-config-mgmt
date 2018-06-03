package main

type K8sConfigManagement struct {
	Configuration Configuration
	K8sService Kubernetes
	Logger *DaemonLogger

	namespaces map[string]cfgNamespace
	clusterConfig cfgCluster
}

func (mgmt *K8sConfigManagement) Run() {
	mgmt.Init()
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

func (mgmt *K8sConfigManagement) ManageConfiguration() {
	mgmt.Logger.Main("Manage Configuration")

	mgmt.Logger.Category("Cluster configuration")
	mgmt.ManagePodSecurityPolicies()
	mgmt.ManageClusterRoles()
	mgmt.ManageClusterRoleBindings()
	mgmt.ManageStorageClasses()
	mgmt.ManageNamespaces()

	for _, namespace := range mgmt.namespaces {
		mgmt.Logger.Category("Namespace %v", namespace.Name)

		mgmt.ManageNamespaceServiceAccounts(namespace)
		mgmt.ManageNamespaceConfigMaps(namespace)
		mgmt.ManageNamespaceRoles(namespace)
		mgmt.ManageNamespaceRoleBindings(namespace)
		mgmt.ManageNamespaceResourceQuotas(namespace)
		mgmt.ManageNamespacePodPresets(namespace)
		mgmt.ManageNamespacePodDisruptionBudgets(namespace)
		mgmt.ManageNamespaceLimitRanges(namespace)
		mgmt.ManageNamespaceNetworkPolicies(namespace)
	}
}
