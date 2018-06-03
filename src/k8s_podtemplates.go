package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/api/core/v1"
)

func (k *Kubernetes) ListPodTemplates(namespace string) (list map[string]v13.PodTemplate, error error) {
	list = map[string]v13.PodTemplate{}

	options := v12.ListOptions{}
	if valList, err := k.Client().CoreV1().PodTemplates(namespace).List(options); err == nil {
		for _, item := range valList.Items {
			list[item.Name] = item
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) CreatePodTemplate(namespace string, PodTemplate *v13.PodTemplate) (psp *v13.PodTemplate, error error) {
	return k.Client().CoreV1().PodTemplates(namespace).Create(PodTemplate)
}

func (k *Kubernetes) UpdatePodTemplate(namespace string, PodTemplate *v13.PodTemplate) (psp *v13.PodTemplate, error error) {
	return k.Client().CoreV1().PodTemplates(namespace).Update(PodTemplate)
}

func (k *Kubernetes) DeletePodTemplate(namespace, name string) (error error) {
	options := v12.DeleteOptions{
		GracePeriodSeconds: &zeroGracePeriod,
	}

	return k.Client().CoreV1().PodTemplates(namespace).Delete(name, &options)
}


