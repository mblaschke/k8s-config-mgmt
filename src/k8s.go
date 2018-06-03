package main

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
	)

const K8S_BLACKLIST_NAMESPACE = "^(kube-system|kube-public|default)"
const K8S_BLACKLIST_SERVICEACCOUNT = "^(default)"
const K8S_BLACKLIST_CLUSTERROLE = "^(admin|cluster-admin|edit|view|system:.+)$"
const K8S_BLACKLIST_CLUSTERROLEBINDING = "^(cluster-admin|minikube-rbac|add-on-cluster-admin|kubeadm:.+|storage-.+|system:.+)$"

type KubernetesBase struct {
	clientset *kubernetes.Clientset

	Logger *DaemonLogger
}

type Kubernetes struct {
	KubernetesBase

	KubeContext string
	KubeConfig string
	AnnotationTrigger string
	AnnotationSelector string
	AnnotationSelectorValue string
}


var (
	k8sBlacklistNamespace = regexp.MustCompile(K8S_BLACKLIST_NAMESPACE)
	k8sBlacklistServiceAccount = regexp.MustCompile(K8S_BLACKLIST_SERVICEACCOUNT)
	k8sBlacklistClusterRole = regexp.MustCompile(K8S_BLACKLIST_CLUSTERROLE)
	k8sBlacklistClusterRoleBinding = regexp.MustCompile(K8S_BLACKLIST_CLUSTERROLEBINDING)
	zeroGracePeriod int64 = 0
)

// Create cached kubernetes client
func (k *Kubernetes) Client() (clientset *kubernetes.Clientset) {
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

func (k *Kubernetes) ClusterRoleBindings() (runtime.Object) {

}
