package main

type K8sConfigManagementBase struct {
	Configuration Configuration
	K8sService Kubernetes
	Logger *DaemonLogger

	namespaces map[string]cfgNamespace
	clusterConfig cfgCluster
}

type K8sConfigManagementBaseCluster struct {
	K8sConfigManagementBase
}

type K8sConfigManagementBaseNamespace struct {
	K8sConfigManagementBase
	Namespace cfgNamespace
}


func (mgmt *K8sConfigManagementBase) filter(name string, whitelist, blacklist []string) (bool) {
	return true
}

func (mgmt *K8sConfigManagementBase) IsNotDryRun() (run bool) {
	if !opts.DryRun {
		run = true
	} else {
		mgmt.Logger.StepResult("dry run")
	}
	return
}

func (mgmt *K8sConfigManagementBase) handleOperationState(err error) {
	if err == nil {
		mgmt.Logger.StepResult("ok")
	} else {
		mgmt.Logger.StepResult("failed [%v]", err)
	}
}

func (mgmt *K8sConfigManagementBase) ClusterScope() (*K8sConfigManagementScopeCluster) {
	return &K8sConfigManagementScopeCluster{
		K8sConfigManagementBaseCluster{K8sConfigManagementBase: *mgmt},
	}
}

func (mgmt *K8sConfigManagementBase) NamespaceScope(namespace cfgNamespace) (*K8sConfigManagementScopeNamespace) {
	return &K8sConfigManagementScopeNamespace{
		K8sConfigManagementBaseNamespace{ K8sConfigManagementBase: *mgmt, Namespace: namespace},
	}
}
