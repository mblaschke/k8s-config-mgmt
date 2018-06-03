package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/extensions/v1beta1"
)

func (k *Kubernetes) ListPodSecurityPolicyies() (list map[string]v1beta1.PodSecurityPolicy, error error) {
	list = map[string]v1beta1.PodSecurityPolicy{}

	options := v12.ListOptions{}

	if valList, err := k.Client().ExtensionsV1beta1().PodSecurityPolicies().List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreatePodSecurityPolicy(PodSecurityPolicy *v1beta1.PodSecurityPolicy) (psp *v1beta1.PodSecurityPolicy, error error) {
	return k.Client().ExtensionsV1beta1().PodSecurityPolicies().Create(PodSecurityPolicy)
}

func (k *Kubernetes) UpdatePodSecurityPolicy(PodSecurityPolicy *v1beta1.PodSecurityPolicy) (psp *v1beta1.PodSecurityPolicy, error error) {
	return k.Client().ExtensionsV1beta1().PodSecurityPolicies().Update(PodSecurityPolicy)
}

func (k *Kubernetes) DeletePodSecurityPolicy(name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().ExtensionsV1beta1().PodSecurityPolicies().Delete(name, &options)
}


