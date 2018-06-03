package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

func (k *Kubernetes) ListResourceQuotas(namespace string) (list map[string]v1.ResourceQuota, error error) {
	list = map[string]v1.ResourceQuota{}

	options := v12.ListOptions{}

	if valList, err := k.Client().CoreV1().ResourceQuotas(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateResourceQuota(namespace string, ResourceQuota *v1.ResourceQuota) (ns *v1.ResourceQuota, error error) {
	return k.Client().CoreV1().ResourceQuotas(namespace).Create(ResourceQuota)
}

func (k *Kubernetes) UpdateResourceQuota(namespace string, ResourceQuota *v1.ResourceQuota) (ns *v1.ResourceQuota, error error) {
	return k.Client().CoreV1().ResourceQuotas(namespace).Update(ResourceQuota)
}

func (k *Kubernetes) DeleteResourceQuota(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().ResourceQuotas(namespace).Delete(name, &options)
}


