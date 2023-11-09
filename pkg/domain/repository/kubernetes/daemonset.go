package kubernetes

import (
	"context"

	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) GetDaemonset(name, namespace string) (*v1.DaemonSet, error) {
	return k8sRepo.k8sClient.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}
