package kubernetes

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubectl "k8s.io/kubectl/pkg/cmd/create"
)

// Создаем Image Pull Secret и при необходимости привязываем его к ServiceAccounts
func (k8sRepo *k8sRepo) CreatePullSecret(ctx context.Context, params params.PullSecretParams) error {
	// todo: retry
	if secret, err := k8sRepo.k8sClient.CoreV1().Secrets(params.Namespace).Get(ctx, params.Name, metav1.GetOptions{}); err != nil {
		if kerrors.IsNotFound(err) {
			// Не нашли секрет, создаем новый
			log.Println(err.Error())
			log.Println("Create pull secret", params.Name)
			// Создаем содержимое секрета
			if data, err := handleDockerCfgJSONContent(params.User, params.Password, params.Server); err != nil {
				return err
			} else {
				newSecret := &corev1.Secret{
					TypeMeta: metav1.TypeMeta{
						APIVersion: corev1.SchemeGroupVersion.String(),
						Kind:       "Secret",
					},
					ObjectMeta: metav1.ObjectMeta{
						Name:      params.Name,
						Namespace: params.Namespace,
					},
					Type: corev1.SecretTypeDockerConfigJson,
					Data: map[string][]byte{
						corev1.DockerConfigJsonKey: data,
					},
				}

				// Создаем сам секрет
				if _, err := k8sRepo.k8sClient.CoreV1().Secrets(params.Namespace).Create(ctx, newSecret, metav1.CreateOptions{}); err != nil {
					return err
				}
			}
		} else {
			return err
		}

	} else {
		// Если секрет есть, но у него другой типа - возвращаем ошибку (todo: у Image Pull Secret несколько типов на самом деле)
		if secret.Type != corev1.SecretTypeDockerConfigJson {
			return errors.New(fmt.Sprintln("Invalid secret type", secret.Type, "should be", corev1.SecretTypeDockerConfigJson))
		}
		// Создаем содержимое секрета
		if data, err := handleDockerCfgJSONContent(params.User, params.Password, params.Server); err != nil {
			return err
		} else {
			// Проверяем есть ли необходимость обновлять секрет
			var update bool
			if dataFromCluster, ok := secret.Data[corev1.DockerConfigJsonKey]; !ok {
				// Если такого ключа нет, то обновляем
				update = true
			} else {
				// Если ключ есть, но его содержимое отличается, то тоже обновляем
				if !reflect.DeepEqual(dataFromCluster, data) {
					update = true
				}
			}

			if update {
				secret.Data[corev1.DockerConfigJsonKey] = data
				// Обновляшка
				if _, err := k8sRepo.k8sClient.CoreV1().Secrets(params.Namespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
					return err
				}
				log.Println("Secret", params.Name, "is uptaded in namespace", params.Namespace)
			}
		}
	}

	// Для каждого указанного ServiceAccount приписываем данный Image Pull Secret
	log.Println("Secret", params.Name, "is already in namespace", params.Namespace)
	for _, sa := range params.ServiceAccounts {
		log.Println("Add pull secret", params.Name, "to sa", sa, "in namespace", params.Namespace)
		if err := k8sRepo.addPullSecretToSa(ctx, params.Name, params.Namespace, sa); err != nil {
			return err
		}
	}

	return nil
}

// FROM: https://github.com/kubernetes/kubectl/blob/197123726db24c61aa0f78d1f0ba6e91a2ec2f35/pkg/cmd/create/create_secret_docker.go

// handleDockerCfgJSONContent serializes a ~/.docker/config.json file
func handleDockerCfgJSONContent(username, password, server string) ([]byte, error) {
	dockerConfigAuth := kubectl.DockerConfigEntry{
		Username: username,
		Password: password,
		//Email:    email,
		Auth: encodeDockerConfigFieldAuth(username, password),
	}
	dockerConfigJSON := kubectl.DockerConfigJSON{
		Auths: map[string]kubectl.DockerConfigEntry{server: dockerConfigAuth},
	}

	return json.Marshal(dockerConfigJSON)
}

// encodeDockerConfigFieldAuth returns base64 encoding of the username and password string
func encodeDockerConfigFieldAuth(username, password string) string {
	fieldValue := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(fieldValue))
}
