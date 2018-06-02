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
	k8sServiceAccountDefaultBlacklist = regexp.MustCompile(K8S_BLACKLIST_SERVICEACCOUNT)
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
