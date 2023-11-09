package kubernetes

import (
	"context"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

// Добавить Image Pull Secret к ServiceAccount
func (k8sRepo *k8sRepo) addPullSecretToSa(ctx context.Context, secret, ns, sa string) error {
	// Проверяем есть ли в целом указанный SA
	return retry.RetryOnConflict(wait.Backoff{
		Steps:    10,
		Duration: 5 * time.Second,
		Factor:   1.5,
		Jitter:   0.1,
	}, func() error {
		if saFromCluster, err := k8sRepo.k8sClient.CoreV1().ServiceAccounts(ns).Get(ctx, sa, metav1.GetOptions{}); err != nil {
			// Если нет, то кидаем ошибку
			if kerrors.IsNotFound(err) {
				// todo: кидать ошибку ?
				log.Println(err.Error())
				return nil
			}
			return err
		} else {
			// Если SA есть, ищем наш секрет среди его ImagePullSecrets
			for _, name := range saFromCluster.ImagePullSecrets {
				// На месте, ничего не делаем
				if name.Name == secret {
					return nil
				}
			}

			// Добавляем нашего шалунишку
			saFromCluster.ImagePullSecrets = append(saFromCluster.ImagePullSecrets, corev1.LocalObjectReference{Name: secret})
			// Обновляем
			if _, err := k8sRepo.k8sClient.CoreV1().ServiceAccounts(ns).Update(ctx, saFromCluster, metav1.UpdateOptions{}); err != nil {
				return err
			}
			return nil
		}
	})
}
