package main

import (
	"regexp"
	"errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

const K8S_BLACKLIST_NAMESPACE = "^(kube-system|kube-public|default)"

type Kubernetes struct {
	clientset *kubernetes.Clientset

	Logger *DaemonLogger

	KubeContext string
	KubeConfig string
	AnnotationTrigger string
	AnnotationSelector string
	AnnotationSelectorValue string
}

var (
	k8sNamespaceBlacklist = regexp.MustCompile(K8S_BLACKLIST_NAMESPACE)
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

func (k *Kubernetes) ListNamespaces() (list map[string]v1.Namespace, error error) {
	list = map[string]v1.Namespace{}

	options := v12.ListOptions{}

	if valList, err := k.Client().CoreV1().Namespaces().List(options); err == nil {
		for _, item := range valList.Items {
			if !k8sNamespaceBlacklist.MatchString(item.Name) {
				list[item.Name] = item
			}
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateNamespace(namespace v1.Namespace) (ns *v1.Namespace, error error) {
	if k8sNamespaceBlacklist.MatchString(namespace.Name) {
		return nil, errors.New("Cannot create blacklisted namespace")
	}

	return k.Client().CoreV1().Namespaces().Create(&namespace)
}

func (k *Kubernetes) UpdateNamespace(namespace v1.Namespace) (ns *v1.Namespace, error error) {
	if k8sNamespaceBlacklist.MatchString(namespace.Name) {
		return nil, errors.New("Cannot update blacklisted namespace")
	}

	return k.Client().CoreV1().Namespaces().Update(&namespace)
}

func (k *Kubernetes) DeleteNamespace(name string) (error error) {
	if k8sNamespaceBlacklist.MatchString(name) {
		return errors.New("Cannot delete blacklisted namespace")
	}

	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().Namespaces().Delete(name, &options)
}

