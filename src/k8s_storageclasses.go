package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/storage/v1"
)

func (k *Kubernetes) ListStorageClasses() (list map[string]v1.StorageClass, error error) {
	list = map[string]v1.StorageClass{}

	options := v12.ListOptions{}

	if valList, err := k.Client().StorageV1().StorageClasses().List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreateStorageClass(StorageClass *v1.StorageClass) (psp *v1.StorageClass, error error) {
	return k.Client().StorageV1().StorageClasses().Create(StorageClass)
}

func (k *Kubernetes) UpdateStorageClass(StorageClass *v1.StorageClass) (psp *v1.StorageClass, error error) {
	return k.Client().StorageV1().StorageClasses().Update(StorageClass)
}

func (k *Kubernetes) DeleteStorageClass(name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().StorageV1().StorageClasses().Delete(name, &options)
}


