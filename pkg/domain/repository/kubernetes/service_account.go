package kubernetes

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Получить ServiceAccount
func (k8sRepo *k8sRepo) GetServiceAccount(name, namespace string) (*v1.ServiceAccount, error) {
	return k8sRepo.k8sClient.CoreV1().ServiceAccounts(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
