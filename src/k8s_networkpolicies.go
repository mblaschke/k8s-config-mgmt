package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/networking/v1"
)

func (k *Kubernetes) ListNetworkPolicies(namespace string) (list map[string]v1.NetworkPolicy, error error) {
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

func (k *Kubernetes) CreateNetworkPolicy(namespace string, NetworkPolicy *v1.NetworkPolicy) (psp *v1.NetworkPolicy, error error) {
	return k.Client().NetworkingV1().NetworkPolicies(namespace).Create(NetworkPolicy)
}

func (k *Kubernetes) UpdateNetworkPolicy(namespace string, NetworkPolicy *v1.NetworkPolicy) (psp *v1.NetworkPolicy, error error) {
	return k.Client().NetworkingV1().NetworkPolicies(namespace).Update(NetworkPolicy)
}

func (k *Kubernetes) DeleteNetworkPolicy(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().NetworkingV1().NetworkPolicies(namespace).Delete(name, &options)
}


