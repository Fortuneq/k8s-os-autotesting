package kubernetes

import (
	"context"
	"log"
	"time"

	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"
	"sbet-tech.com/synapse/istio-se/allure/pkg/utils"
)

func (k8sRepo *k8sRepo) CreateTestAppDeployment(ctx context.Context, params params.TestAppParams) error {
	// Есть ли неймспейс для тестового приложения
	if err := k8sRepo.CheckNamespace(params.Namespace); err != nil {
		return err
	}

	// Деплоим деплоймент
	if err := k8sRepo.DeployDeployment(ctx, params.AppName, params.Namespace, params.Image, params.Labels, params.Annotations); err != nil {
		return err
	}

	// Получаем имя пода, чтобы проверить, что сайдкар подтянулся
	podName, err := (utils.RetryWithResponse{
		N:     10,
		Sleep: 10 * time.Second,
		Fn: func() (string, error) {
			return k8sRepo.GetPodName(ctx, params.AppName, params.Namespace)
		},
	}).Start()
	if err != nil {
		return err
	}

	// Readiness probe для Istio контейнера
	_, err = (utils.RetryWithResponse{
		N:     10,
		Sleep: 10 * time.Second,
		Fn: func() (string, error) {
			return k8sRepo.PodSendRequest(podName, params.Namespace, "http://0.0.0.0:15021/healthz/ready", "15021:15021")
		},
	}).Start()

	return err
}

func (k8sRepo *k8sRepo) CreateTestAppPod(ctx context.Context, params params.TestAppParams) error {
	// Есть ли неймспейс для тестового приложения
	if err := k8sRepo.CheckNamespace(params.Namespace); err != nil {
		return err
	}

	// Разворачиваем под
	if err := k8sRepo.DeployPod(ctx, params.AppName, params.Namespace, params.Image, params.Labels, params.Annotations); err != nil {
		return err
	}

	// Readiness probe для Istio контейнера
	_, err := (utils.RetryWithResponse{
		N:     10,
		Sleep: 10 * time.Second,
		Fn: func() (string, error) {
			return k8sRepo.PodSendRequest(params.AppName, params.Namespace, "http://0.0.0.0:15021/healthz/ready", "15021:15021")
		},
	}).Start()

	return err
}

func (k8sRepo *k8sRepo) DeleteTestAppDeployment(ctx context.Context, params params.TestAppParams) error {
	if err := k8sRepo.DeleteDeployment(ctx, params.AppName, params.Namespace); err != nil {
		if kerrors.IsNotFound(err) {
			log.Println(err.Error())
		} else {
			return err
		}
	}
	log.Printf("Test app deployment %s was deleted from namespace %s", params.AppName, params.Namespace)
	return nil
}

func (k8sRepo *k8sRepo) DeleteTestAppPod(ctx context.Context, params params.TestAppParams) error {
	if err := k8sRepo.DeletePod(ctx, params.AppName, params.Namespace); err != nil {
		if kerrors.IsNotFound(err) {
			log.Println(err.Error())
		} else {
			return err
		}
	}
	log.Printf("Test app pod %s was deleted from namespace %s", params.AppName, params.Namespace)
	return nil
}
