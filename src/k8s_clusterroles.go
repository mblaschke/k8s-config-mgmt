package main

import (
	"errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/api/rbac/v1"
	)

type KubernetesServiceClusterRoles struct {
	KubernetesBase
}

func (k *KubernetesServiceClusterRoles) List() (list map[string]v13.ClusterRole, error error) {
	list = map[string]v13.ClusterRole{}

	options := v12.ListOptions{}

	if valList, err := k.Client().RbacV1().ClusterRoles().List(options); err == nil {
		for _, item := range valList.Items {
			// disable rbac defaults
			if _, ok := item.Labels["kubernetes.io/bootstrapping"]; ok {
				continue
			}

			if !k8sBlacklistClusterRole.MatchString(item.Name) {
				list[item.Name] = item
			}
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServiceClusterRoles) Create(clusterRole *v13.ClusterRole) (ns *v13.ClusterRole, error error) {
	if k8sBlacklistClusterRole.MatchString(clusterRole.Name) {
		return nil, errors.New("Cannot create blacklisted ClusterRole")
	}

	return k.Client().RbacV1().ClusterRoles().Create(clusterRole)
}

func (k *KubernetesServiceClusterRoles) Update(clusterRole *v13.ClusterRole) (ns *v13.ClusterRole, error error) {
	if k8sBlacklistClusterRole.MatchString(clusterRole.Name) {
		return nil, errors.New("Cannot update blacklisted ClusterRole")
	}

	return k.Client().RbacV1().ClusterRoles().Update(clusterRole)
}

func (k *KubernetesServiceClusterRoles) Delete(name string) (error error) {
	if k8sBlacklistClusterRole.MatchString(name) {
		return errors.New("Cannot delete blacklisted ClusterRole")
	}

	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().RbacV1().ClusterRoles().Delete(name, &options)
}


