package configmanagement

import (
	"k8s-config-mgmt/src/config"
	"k8s-config-mgmt/src/k8s"
	"k8s-config-mgmt/src/logger"
)

type K8sConfigManagementBase struct {
	Configuration config.Configuration
	K8sService k8s.Kubernetes
	Logger *logger.DaemonLogger

	namespaces map[string]config.ConfigNamespace
	clusterConfig config.ConfigCluster
	DryRun bool
	Validate bool
}

func (mgmt *K8sConfigManagementBase) filter(name string, whitelist, blacklist []string) (bool) {
	return true
}

func (mgmt *K8sConfigManagementBase) IsNotDryRun() (run bool) {
	if !mgmt.DryRun {
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

func (mgmt *K8sConfigManagementBase) NamespaceScope(namespace config.ConfigNamespace) (*K8sConfigManagementScopeNamespace) {
	return &K8sConfigManagementScopeNamespace{
		K8sConfigManagementBaseNamespace{ K8sConfigManagementBase: *mgmt, Namespace: namespace},
	}
}
