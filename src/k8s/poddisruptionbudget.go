package k8s

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/policy/v1beta1"
)

type KubernetesServicePodDisruptionBudgets struct {
	KubernetesBase
}

func (k *KubernetesServicePodDisruptionBudgets) List(namespace string) (list map[string]v1beta1.PodDisruptionBudget, error error) {
	list = map[string]v1beta1.PodDisruptionBudget{}

	options := v12.ListOptions{}
	if valList, err := k.Client().PolicyV1beta1().PodDisruptionBudgets(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServicePodDisruptionBudgets) Create(namespace string, PodDisruptionBudget *v1beta1.PodDisruptionBudget) (psp *v1beta1.PodDisruptionBudget, error error) {
	return k.Client().PolicyV1beta1().PodDisruptionBudgets(namespace).Create(PodDisruptionBudget)
}

func (k *KubernetesServicePodDisruptionBudgets) Update(namespace string, PodDisruptionBudget *v1beta1.PodDisruptionBudget) (psp *v1beta1.PodDisruptionBudget, error error) {
	return k.Client().PolicyV1beta1().PodDisruptionBudgets(namespace).Update(PodDisruptionBudget)
}

func (k *KubernetesServicePodDisruptionBudgets) Delete(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().PolicyV1beta1().PodDisruptionBudgets(namespace).Delete(name, &options)
}


