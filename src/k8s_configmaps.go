package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

func (k *Kubernetes) ListConfigMaps(namespace string) (list map[string]v1.ConfigMap, error error) {
	list = map[string]v1.ConfigMap{}

	options := v12.ListOptions{}

	if valList, err := k.Client().CoreV1().ConfigMaps(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateConfigMap(namespace string, configMap *v1.ConfigMap) (ns *v1.ConfigMap, error error) {
	return k.Client().CoreV1().ConfigMaps(namespace).Create(configMap)
}

func (k *Kubernetes) UpdateConfigMap(namespace string, configMap *v1.ConfigMap) (ns *v1.ConfigMap, error error) {
	return k.Client().CoreV1().ConfigMaps(namespace).Update(configMap)
}

func (k *Kubernetes) DeleteConfigMap(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().ConfigMaps(namespace).Delete(name, &options)
}


