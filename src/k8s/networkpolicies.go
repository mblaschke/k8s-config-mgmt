package k8s

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/networking/v1"
)

type KubernetesServiceNetworkPolicies struct {
	KubernetesBase
}

func (k *KubernetesServiceNetworkPolicies) List(namespace string) (list map[string]v1.NetworkPolicy, error error) {
	list = map[string]v1.NetworkPolicy{}

	options := v12.ListOptions{}
	
	if valList, err := k.Client().NetworkingV1().NetworkPolicies(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServiceNetworkPolicies) Create(namespace string, NetworkPolicy *v1.NetworkPolicy) (psp *v1.NetworkPolicy, error error) {
	return k.Client().NetworkingV1().NetworkPolicies(namespace).Create(NetworkPolicy)
}

func (k *KubernetesServiceNetworkPolicies) Update(namespace string, NetworkPolicy *v1.NetworkPolicy) (psp *v1.NetworkPolicy, error error) {
	return k.Client().NetworkingV1().NetworkPolicies(namespace).Update(NetworkPolicy)
}

func (k *KubernetesServiceNetworkPolicies) Delete(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().NetworkingV1().NetworkPolicies(namespace).Delete(name, &options)
}


