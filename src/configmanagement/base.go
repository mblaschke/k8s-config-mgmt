package configmanagement

import (
	"k8s-config-mgmt/src/config"
	"k8s-config-mgmt/src/k8s"
	"k8s-config-mgmt/src/logger"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type K8sConfigManagementBase struct {
	GlobalConfiguration config.Configuration
	Configuration config.ConfigurationManagementItem
	K8sService *k8s.Kubernetes
	Logger *logger.DaemonLogger

	namespaces map[string]config.ConfigNamespace
	clusterConfig config.ConfigCluster
	DryRun bool
	Validate bool
	Force bool
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
		panic(err)
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

func (mgmt *K8sConfigManagementBase) isK8sObjectFiltered(k8sObject runtime.Object) (ret bool) {
	ret = false

	if ! mgmt.isK8sObjectWhitelisted(k8sObject) {
		ret = true
	}

	if mgmt.isK8sObjectBlacklisted(k8sObject) {
		ret = true
	}

	return
}

func (mgmt *K8sConfigManagementBase) isK8sObjectWhitelisted(k8sObject runtime.Object) (ret bool) {
	ret = false

	// whitelist
	if mgmt.Configuration.Whitelist != nil {
		for _, whitelist := range mgmt.Configuration.Whitelist {
			regexp := filterValueToRegexp(*whitelist)
			if regexp.MatchString(k8sObject.(v1.Object).GetName()) {
				ret = true
				break
			}
		}
	} else {
		ret = true
	}

	return
}


func (mgmt *K8sConfigManagementBase) isK8sObjectBlacklisted(k8sObject runtime.Object) (ret bool) {
	ret = false

	// blacklist
	if mgmt.Configuration.Blacklist != nil {
		for _, blacklist := range mgmt.Configuration.Blacklist {
			regexp := filterValueToRegexp(*blacklist)
			if regexp.MatchString(k8sObject.(v1.Object).GetName()) {
				ret = true
				break
			}
		}
	}

	return
}
