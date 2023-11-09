package kubernetes

import (
	"context"

	rbacV1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) GetRole(name, namespace string) (*rbacV1.Role, error) {
	return k8sRepo.k8sClient.RbacV1().Roles(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
