package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) GetClusterRole(name string) (*v1.ClusterRole, error) {
	return k8sRepo.k8sClient.RbacV1().ClusterRoles().Get(context.TODO(), name, metav1.GetOptions{})
}
