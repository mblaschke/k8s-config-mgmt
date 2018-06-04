package k8s

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta12 "k8s.io/api/policy/v1beta1"
)

type KubernetesServicePodSecurityPolicies struct {
	KubernetesBase
}

func (k *KubernetesServicePodSecurityPolicies) List() (list map[string]v1beta12.PodSecurityPolicy, error error) {
	list = map[string]v1beta12.PodSecurityPolicy{}

	options := v12.ListOptions{}

	if valList, err := k.Client().PolicyV1beta1().PodSecurityPolicies().List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServicePodSecurityPolicies) Create(PodSecurityPolicy *v1beta12.PodSecurityPolicy) (psp *v1beta12.PodSecurityPolicy, error error) {
	return k.Client().PolicyV1beta1().PodSecurityPolicies().Create(PodSecurityPolicy)
}

func (k *KubernetesServicePodSecurityPolicies) Update(PodSecurityPolicy *v1beta12.PodSecurityPolicy) (psp *v1beta12.PodSecurityPolicy, error error) {
	return k.Client().PolicyV1beta1().PodSecurityPolicies().Update(PodSecurityPolicy)
}

func (k *KubernetesServicePodSecurityPolicies) Delete(name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().ExtensionsV1beta1().PodSecurityPolicies().Delete(name, &options)
}


