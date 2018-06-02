package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

func (k *Kubernetes) ListLimitRanges(namespace string) (list map[string]v1.LimitRange, error error) {
	list = map[string]v1.LimitRange{}

	options := v12.ListOptions{}

	if valList, err := k.Client().CoreV1().LimitRanges(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateLimitRange(namespace string, limitRange *v1.LimitRange) (ns *v1.LimitRange, error error) {
	return k.Client().CoreV1().LimitRanges(namespace).Create(limitRange)
}

func (k *Kubernetes) UpdateLimitRange(namespace string, limitRange *v1.LimitRange) (ns *v1.LimitRange, error error) {
	return k.Client().CoreV1().LimitRanges(namespace).Update(limitRange)
}

func (k *Kubernetes) DeleteLimitRange(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().LimitRanges(namespace).Delete(name, &options)
}


