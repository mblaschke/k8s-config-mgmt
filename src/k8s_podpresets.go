package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1alpha1 "k8s.io/api/settings/v1alpha1"
)

func (k *Kubernetes) ListPodPresets(namespace string) (list map[string]v1alpha1.PodPreset, error error) {
	list = map[string]v1alpha1.PodPreset{}

	options := v12.ListOptions{}
	
	if valList, err := k.Client().SettingsV1alpha1().PodPresets(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreatePodPreset(namespace string, PodPreset *v1alpha1.PodPreset) (psp *v1alpha1.PodPreset, error error) {
	return k.Client().SettingsV1alpha1().PodPresets(namespace).Create(PodPreset)
}

func (k *Kubernetes) UpdatePodPreset(namespace string, PodPreset *v1alpha1.PodPreset) (psp *v1alpha1.PodPreset, error error) {
	return k.Client().SettingsV1alpha1().PodPresets(namespace).Update(PodPreset)
}

func (k *Kubernetes) DeletePodPreset(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().SettingsV1alpha1().PodPresets(namespace).Delete(name, &options)
}


