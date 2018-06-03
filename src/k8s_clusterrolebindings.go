package main

import (
	"errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/api/rbac/v1"
	)

func (k *Kubernetes) ListClusterRoleBindings() (list map[string]v13.ClusterRoleBinding, error error) {
	list = map[string]v13.ClusterRoleBinding{}

	options := v12.ListOptions{}

	if valList, err := k.Client().RbacV1().ClusterRoleBindings().List(options); err == nil {
		for _, item := range valList.Items {
			// disable rbac defaults
			if _, ok := item.Labels["kubernetes.io/bootstrapping"]; ok {
				continue
			}

			if !k8sBlacklistClusterRoleBinding.MatchString(item.Name) {
				list[item.Name] = item
			}
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateClusterRoleBinding(ClusterRoleBinding *v13.ClusterRoleBinding) (ns *v13.ClusterRoleBinding, error error) {
	if k8sBlacklistClusterRoleBinding.MatchString(ClusterRoleBinding.Name) {
		return nil, errors.New("Cannot create blacklisted ClusterRoleBinding")
	}

	return k.Client().RbacV1().ClusterRoleBindings().Create(ClusterRoleBinding)
}

func (k *Kubernetes) UpdateClusterRoleBinding(ClusterRoleBinding *v13.ClusterRoleBinding) (ns *v13.ClusterRoleBinding, error error) {
	if k8sBlacklistClusterRoleBinding.MatchString(ClusterRoleBinding.Name) {
		return nil, errors.New("Cannot update blacklisted ClusterRoleBinding")
	}

	return k.Client().RbacV1().ClusterRoleBindings().Update(ClusterRoleBinding)
}

func (k *Kubernetes) DeleteClusterRoleBinding(name string) (error error) {
	if k8sBlacklistClusterRoleBinding.MatchString(name) {
		return errors.New("Cannot delete blacklisted ClusterRoleBinding")
	}

	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().RbacV1().ClusterRoleBindings().Delete(name, &options)
}


