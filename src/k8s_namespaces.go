package main

import (
	"errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	)

func (k *Kubernetes) ListNamespaces() (list map[string]*v1.Namespace, error error) {
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

func (k *Kubernetes) CreateNamespace(namespace *v1.Namespace) (ns *v1.Namespace, error error) {
	if k8sBlacklistNamespace.MatchString(namespace.Name) {
		return nil, errors.New("Cannot create blacklisted namespace")
	}

	return k.Client().CoreV1().Namespaces().Create(namespace)
}

func (k *Kubernetes) UpdateNamespace(namespace *v1.Namespace) (ns *v1.Namespace, error error) {
	if k8sBlacklistNamespace.MatchString(namespace.Name) {
		return nil, errors.New("Cannot update blacklisted namespace")
	}

	return k.Client().CoreV1().Namespaces().Update(namespace)
}

func (k *Kubernetes) DeleteNamespace(name string) (error error) {
	if k8sBlacklistNamespace.MatchString(name) {
		return errors.New("Cannot delete blacklisted namespace")
	}

	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().Namespaces().Delete(name, &options)
}


