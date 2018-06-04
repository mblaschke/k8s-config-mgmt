package k8s

import (
	"regexp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/kubernetes/scheme"
	"io/ioutil"
	"fmt"
	"log"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s-config-mgmt/src/logger"
)

const K8S_BLACKLIST_NAMESPACE = "^(kube-system|kube-public|default)"
const K8S_BLACKLIST_SERVICEACCOUNT = "^(default)"
const K8S_BLACKLIST_CLUSTERROLE = "^(admin|cluster-admin|edit|view|system:.+)$"
const K8S_BLACKLIST_CLUSTERROLEBINDING = "^(cluster-admin|minikube-rbac|add-on-cluster-admin|kubeadm:.+|storage-.+|system:.+)$"

type KubernetesBase struct {
	clientset *kubernetes.Clientset

	Logger *logger.DaemonLogger
	KubeContext string
	KubeConfig string
}

type Kubernetes struct {
	KubernetesBase
}


var (
	k8sBlacklistNamespace = regexp.MustCompile(K8S_BLACKLIST_NAMESPACE)
	k8sBlacklistServiceAccount = regexp.MustCompile(K8S_BLACKLIST_SERVICEACCOUNT)
	k8sBlacklistClusterRole = regexp.MustCompile(K8S_BLACKLIST_CLUSTERROLE)
	k8sBlacklistClusterRoleBinding = regexp.MustCompile(K8S_BLACKLIST_CLUSTERROLEBINDING)
	zeroGracePeriod int64 = 0
)

// Create cached kubernetes client
func (k *KubernetesBase) Client() (clientset *kubernetes.Clientset) {
	var err error
	var config *rest.Config

	if k.clientset == nil {
		if k.KubeConfig != "" {
			// KUBECONFIG
			config, err = buildConfigFromFlags(k.KubeContext, k.KubeConfig)
			if err != nil {
				panic(err.Error())
			}
		} else {
			// K8S in cluster
			config, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}

		k.clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	return k.clientset
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func (k *Kubernetes) ParseConfig(path string) (runtime.Object) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(data, nil, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
	}
	return obj
}

func (k *Kubernetes) ClusterRoleBindings() (*KubernetesServiceClusterRoleBindings) {
	return &KubernetesServiceClusterRoleBindings{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) ClusterRoles() (*KubernetesServiceClusterRoles) {
	return &KubernetesServiceClusterRoles{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) ConfigMaps() (*KubernetesServiceConfigMaps) {
	return &KubernetesServiceConfigMaps{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) LimitRanges() (*KubernetesServiceLimitRanges) {
	return &KubernetesServiceLimitRanges{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) Namespaces() (*KubernetesServiceNamespaces) {
	return &KubernetesServiceNamespaces{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) NetworkPolicies() (*KubernetesServiceNetworkPolicies) {
	return &KubernetesServiceNetworkPolicies{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) PodDisruptionBudgets() (*KubernetesServicePodDisruptionBudgets) {
	return &KubernetesServicePodDisruptionBudgets{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) PodPresets() (*KubernetesServicePodPresets) {
	return &KubernetesServicePodPresets{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) PodSecurityPolicies() (*KubernetesServicePodSecurityPolicies) {
	return &KubernetesServicePodSecurityPolicies{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) PodTemplates() (*KubernetesServicePodTemplates) {
	return &KubernetesServicePodTemplates{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) ResourceQuotas() (*KubernetesServiceResourceQuotas) {
	return &KubernetesServiceResourceQuotas{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) RoleBindings() (*KubernetesServiceRoleBindings) {
	return &KubernetesServiceRoleBindings{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) Roles() (*KubernetesServiceRoles) {
	return &KubernetesServiceRoles{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) ServiceAccounts() (*KubernetesServiceServiceAccounts) {
	return &KubernetesServiceServiceAccounts{KubernetesBase: k.KubernetesBase}
}

func (k *Kubernetes) StorageClasses() (*KubernetesServiceStorageClasses) {
	return &KubernetesServiceStorageClasses{KubernetesBase: k.KubernetesBase}
}
