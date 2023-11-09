package kubernetes

import (
	"context"
	"log"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"
	"sbet-tech.com/synapse/istio-se/allure/pkg/utils"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) CreateNamespace(ctx context.Context, params params.NamespaceParams) error {
	// todo: retry
	if nsFromCluster, err := k8sRepo.k8sClient.CoreV1().Namespaces().Get(ctx, params.Name, metav1.GetOptions{}); err != nil {
		// Уже есть, чекаем лейблы
		if kerrors.IsNotFound(err) {
			log.Println(err.Error())

			log.Println("Create namespace", params.Name)
			ns := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
				Name:   params.Name,
				Labels: params.Labels,
			}}

			_, err = k8sRepo.k8sClient.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})

		}

		return err
	} else {
		// Лейблы не сходятся, нужно обновиться
		if !utils.CompareMaps(params.Labels, nsFromCluster.Labels) {
			nsFromCluster.DeepCopy().SetLabels(params.Labels)

			log.Println("Update namespace labels", params.Name)
			if _, err := k8sRepo.k8sClient.CoreV1().Namespaces().Update(ctx, nsFromCluster, metav1.UpdateOptions{}); err != nil {
				return err
			}

			return nil
		}

	}

	if len(params.Labels) == 0 {
		log.Println("Namespace", params.Name, "is already in cluster with no labels")
	} else {
		log.Println("Namespace", params.Name, "is already in cluster with required labels", params.Labels)
	}

	return nil
}

func (k8sRepo *k8sRepo) CheckNamespace(namespace string) error {
	_, err := k8sRepo.k8sClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (k8sRepo *k8sRepo) DeleteNamespace(ctx context.Context, params params.NamespaceParams) error {

	// todo: retry
	if err := k8sRepo.k8sClient.CoreV1().Namespaces().Delete(ctx, params.Name, metav1.DeleteOptions{}); err != nil {
		// Уже есть, чекаем лейблы
		if kerrors.IsNotFound(err) {
			log.Println(err.Error())
		} else {
			return err
		}

	}

	log.Println("Namespace", params.Name, "was deleted")

	return nil
}
