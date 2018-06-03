package main

import (
	"errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	)

type KubernetesServiceNamespaces struct {
	KubernetesBase
}

func (k *KubernetesServiceNamespaces) List() (list map[string]*v1.Namespace, error error) {
	list = map[string]*v1.Namespace{}

	options := v12.ListOptions{}

	if valList, err := k.Client().CoreV1().Namespaces().List(options); err == nil {
		for key, item := range valList.Items {
			if !k8sBlacklistNamespace.MatchString(item.Name) {
				list[item.Name] = &valList.Items[key]
			}
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServiceNamespaces) Create(namespace *v1.Namespace) (ns *v1.Namespace, error error) {
	if k8sBlacklistNamespace.MatchString(namespace.Name) {
		return nil, errors.New("Cannot create blacklisted namespace")
	}

	return k.Client().CoreV1().Namespaces().Create(namespace)
}

func (k *KubernetesServiceNamespaces) Update(namespace *v1.Namespace) (ns *v1.Namespace, error error) {
	if k8sBlacklistNamespace.MatchString(namespace.Name) {
		return nil, errors.New("Cannot update blacklisted namespace")
	}

	return k.Client().CoreV1().Namespaces().Update(namespace)
}

func (k *KubernetesServiceNamespaces) Delete(name string) (error error) {
	if k8sBlacklistNamespace.MatchString(name) {
		return errors.New("Cannot delete blacklisted namespace")
	}

	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().Namespaces().Delete(name, &options)
}


