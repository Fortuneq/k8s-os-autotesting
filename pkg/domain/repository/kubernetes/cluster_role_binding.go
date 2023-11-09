package kubernetes

import (
	"context"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) GetClusterRoleBinding(name string) (*v1.ClusterRoleBinding, error) {
	return k8sRepo.k8sClient.RbacV1().ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
}
