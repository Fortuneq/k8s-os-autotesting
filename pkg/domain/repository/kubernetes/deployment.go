package kubernetes

import (
	"context"
	"errors"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k8sRepo *k8sRepo) GetDeployment(ctx context.Context, name, namespace string) (*appsv1.Deployment, error) {
	return k8sRepo.k8sClient.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetDeploymentName возвращает под приложения с селектором selector в неймспейсе namespace
func (k8sRepo *k8sRepo) GetDeploymentName(namespace string, selector string) (string, error) {
	ls, err := k8sRepo.k8sClient.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return "", err
	}

	if len(ls.Items) == 0 {
		return "", errors.New("deployment not found")
	}

	dep := ls.Items[0]
	return dep.Name, nil
}

func (k8sRepo *k8sRepo) DeployDeployment(ctx context.Context, appName, namespace, image string, labels, annotations map[string]string) error {
	dep, err := k8sRepo.GetDeployment(ctx, appName, namespace)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Не нашли деплоймент, создаем новый
			if err = k8sRepo.CreateDeployment(ctx, appName, namespace, image, labels, annotations); err != nil {
				return err
			}
			time.Sleep(10 * time.Second)
		} else {
			return err
		}
	} else {
		for _, c := range dep.Spec.Template.Spec.Containers {
			if c.Name == appName {
				c.Image = image
			}
		}
		dep.Spec.Template.Labels = labels
		dep.Spec.Template.Annotations = annotations
		dep.Spec.Selector.MatchLabels = labels
		_, err = k8sRepo.UpdateDeployment(ctx, namespace, dep)
	}
	return err
}

// CreateDeployment создает деплоймент
func (k8sRepo *k8sRepo) CreateDeployment(ctx context.Context, name, ns, image string, labels, annotations map[string]string) error {
	rep := int32(1)
	depl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &rep,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        name,
					Labels:      labels,
					Annotations: annotations,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  name,
							Image: image,
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:    resource.MustParse("200m"),
									v1.ResourceMemory: resource.MustParse("200Mi"),
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := k8sRepo.k8sClient.AppsV1().Deployments(ns).Create(ctx, depl, metav1.CreateOptions{})

	return err
}

// UpdateDeployment обновляет деплоймент
func (k8sRepo *k8sRepo) UpdateDeployment(ctx context.Context, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return k8sRepo.k8sClient.AppsV1().Deployments(namespace).Update(ctx, deployment, metav1.UpdateOptions{})
}

// DeleteDeployment удаляет деплоймент
func (k8sRepo *k8sRepo) DeleteDeployment(ctx context.Context, name, ns string) error {
	return k8sRepo.k8sClient.AppsV1().Deployments(ns).Delete(ctx, name, metav1.DeleteOptions{})
}
