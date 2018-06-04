package k8s

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/api/rbac/v1"
)

type KubernetesServiceRoles struct {
	KubernetesBase
}

func (k *KubernetesServiceRoles) List(namespace string) (list map[string]v13.Role, error error) {
	list = map[string]v13.Role{}

	options := v12.ListOptions{}

	if valList, err := k.Client().RbacV1().Roles(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServiceRoles) Create(namespace string, serviceAccount *v13.Role) (ns *v13.Role, error error) {
	return k.Client().RbacV1().Roles(namespace).Create(serviceAccount)
}

func (k *KubernetesServiceRoles) Update(namespace string, serviceAccount *v13.Role) (ns *v13.Role, error error) {
	return k.Client().RbacV1().Roles(namespace).Update(serviceAccount)
}

func (k *KubernetesServiceRoles) Delete(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().RbacV1().Roles(namespace).Delete(name, &options)
}


