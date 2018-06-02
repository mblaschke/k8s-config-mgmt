package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/api/rbac/v1"
)

func (k *Kubernetes) ListRoleBindings(namespace string) (list map[string]v13.RoleBinding, error error) {
	list = map[string]v13.RoleBinding{}

	options := v12.ListOptions{}

	if valList, err := k.Client().RbacV1().RoleBindings(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateRoleBinding(namespace string, serviceAccount *v13.RoleBinding) (ns *v13.RoleBinding, error error) {
	return k.Client().RbacV1().RoleBindings(namespace).Create(serviceAccount)
}

func (k *Kubernetes) UpdateRoleBinding(namespace string, serviceAccount *v13.RoleBinding) (ns *v13.RoleBinding, error error) {
	return k.Client().RbacV1().RoleBindings(namespace).Update(serviceAccount)
}

func (k *Kubernetes) DeleteRoleBinding(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().RbacV1().RoleBindings(namespace).Delete(name, &options)
}


