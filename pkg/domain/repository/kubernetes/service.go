package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) GetService(name, namespace string) (*v1.Service, error) {
	return k8sRepo.k8sClient.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
