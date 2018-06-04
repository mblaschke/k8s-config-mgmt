package k8s

import (
	"errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

type KubernetesServiceServiceAccounts struct {
	KubernetesBase
}

func (k *KubernetesServiceServiceAccounts) List(namespace string) (list map[string]v1.ServiceAccount, error error) {
	list = map[string]v1.ServiceAccount{}

	options := v12.ListOptions{}

	if valList, err := k.Client().CoreV1().ServiceAccounts(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			if !k8sBlacklistServiceAccount.MatchString(item.Name) {
				list[item.Name] = item
			}
		}
	} else {
		error = err
	}

	return
}

func (k *KubernetesServiceServiceAccounts) Create(namespace string, serviceAccount *v1.ServiceAccount) (ns *v1.ServiceAccount, error error) {
	return k.Client().CoreV1().ServiceAccounts(namespace).Create(serviceAccount)
}

func (k *KubernetesServiceServiceAccounts) Update(namespace string, serviceAccount *v1.ServiceAccount) (ns *v1.ServiceAccount, error error) {
	return k.Client().CoreV1().ServiceAccounts(namespace).Update(serviceAccount)
}

func (k *KubernetesServiceServiceAccounts) Delete(namespace, name string) (error error) {
	if k8sBlacklistServiceAccount.MatchString(name) {
		return errors.New("Cannot delete blacklisted ServiceAccount")
	}

	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().ServiceAccounts(namespace).Delete(name, &options)
}


