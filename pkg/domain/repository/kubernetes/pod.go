package kubernetes

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/model"
	"sbet-tech.com/synapse/istio-se/allure/pkg/utils"
)

type Data struct {
	JsonData []byte
	ApiKey   string
}

type CpuMemory struct {
	Cpu    float64
	Memory float64
}

type CpuMemoryPodsCount struct {
	CpuMemory
	PodsCount int
}

func (k8sRepo k8sRepo) GetCurrentUsageMetrics(selector string, namespace string) (result *CpuMemoryPodsCount) {
	log.Printf("RUN getCurrentUsageMetrics for %s selector", selector)
	podMetricsList, err := k8sRepo.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("getCurrentUsageMetrics get error: %s", err.Error())
	}
	m := &CpuMemoryPodsCount{}
	m.PodsCount = len(podMetricsList.Items)
	for _, podMetrics := range podMetricsList.Items {
		for _, containerMetrics := range podMetrics.Containers {
			cpuQuantity, err := strconv.ParseFloat(containerMetrics.Usage.Cpu().AsDec().String(), 64)
			if err == nil {
				m.Cpu += cpuQuantity
			}
			memoryQuantity, err := strconv.ParseFloat(containerMetrics.Usage.Memory().AsDec().String(), 64)
			if err == nil {
				m.Memory += memoryQuantity
			}
		}
	}
	if m.PodsCount != 0 {

		m.Cpu = m.Cpu / float64(m.PodsCount)
		m.Memory = m.Memory / float64(m.PodsCount)
	}
	log.Printf("getCurrentUsageMetrics USAGE: %v", m)

	return m
}

func (k8sRepo k8sRepo) GetLimits(selector string, namespace string) (result *CpuMemory) {
	log.Printf("RUN getLimits for %s selector", selector)

	podsList, err := k8sRepo.k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Printf("getLimits get error: %s", err.Error())
		fmt.Println(fmt.Errorf("Couldn't get current pods list: %w", err))
	}
	if len(podsList.Items) == 0 {
		return nil
	}
	pod := podsList.Items[0]
	limits := &CpuMemory{}
	for _, container := range pod.Spec.Containers {
		cpuQuantity, err := strconv.ParseFloat(container.Resources.Limits.Cpu().AsDec().String(), 64)
		if err == nil {
			limits.Cpu += cpuQuantity
		}
		memoryQuantity, err := strconv.ParseFloat(container.Resources.Limits.Memory().AsDec().String(), 64)
		if err == nil {
			limits.Memory += memoryQuantity
		}

	}
	log.Printf("getLimits LIMITS : %v", limits)
	return limits
}

// DeployPod Создаем под с приложением
func (k8sRepo k8sRepo) DeployPod(ctx context.Context, name, namespace, image string, labels, annotations map[string]string) error {
	pod, err := k8sRepo.GetPod(ctx, name, namespace)
	if err != nil {
		if kerrors.IsNotFound(err) {
			// Не нашли под, создаем новый
			if err = k8sRepo.CreatePod(ctx, name, namespace, image, labels, annotations); err != nil {
				return err
			}
			time.Sleep(10 * time.Second)
		} else {
			return err
		}
	} else {
		for _, c := range pod.Spec.Containers {
			if c.Name == name {
				c.Image = image
			}
		}
		pod.Labels = labels
		pod.Annotations = annotations
		_, err = k8sRepo.UpdatePod(ctx, namespace, pod)
	}
	return err
}

func (k8sRepo *k8sRepo) GetPod(ctx context.Context, name, namespace string) (*v1.Pod, error) {
	return k8sRepo.k8sClient.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}

func (k8sRepo *k8sRepo) GetPods(ctx context.Context, namespace string) (*v1.PodList, error) {
	return k8sRepo.k8sClient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
}

func (k8sRepo k8sRepo) GetAllDeploymentPodNames(ctx context.Context, deployment, namespace string) ([]string, error) {
	ls, err := k8sRepo.k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	if len(ls.Items) == 0 {
		return nil, kerrors.NewNotFound(v1.Resource("pods"), "pod")
	}

	dep, err := k8sRepo.GetDeployment(ctx, deployment, namespace)
	if err != nil {
		return nil, err
	}
	selector, err := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	if err != nil {
		return nil, err
	}
	replicaSets, err := k8sRepo.k8sClient.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	var result []string
	revisionAnnotation := "deployment.kubernetes.io/revision"
	for _, p := range ls.Items {
		for _, owner := range p.OwnerReferences {
			for _, replicaSet := range replicaSets.Items {
				if owner.Name == replicaSet.Name && dep.Annotations[revisionAnnotation] == replicaSet.Annotations[revisionAnnotation] {
					result = append(result, p.Name)
				}
			}
		}
	}
	return result, nil
}

// GetPodName возвращает под приложения с селектором allure-app=istio-allure в неймспейсе ns // todo комменты поправить
func (k8sRepo *k8sRepo) GetPodName(ctx context.Context, deployment, namespace string) (string, error) {
	// Получаем деплоймент по имени
	dep, err := k8sRepo.GetDeployment(ctx, deployment, namespace)
	if err != nil {
		return "", err
	}

	// Получаем его селекторы
	selector, err := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	if err != nil {
		return "", err
	}

	// У деплоймента всегда есть эта аннотация
	revision := dep.Annotations["deployment.kubernetes.io/revision"]

	// Ищем реплика сеты по этой аннотации
	replicaSets, err := k8sRepo.k8sClient.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: selector.String(),
	})
	if err != nil {
		return "", err
	}

	if len(replicaSets.Items) == 0 {
		return "", kerrors.NewNotFound(v1.Resource("replicaSets"), "replicaSet")
	}

	// Ищем репликасеты, которые были созданы для этого деплоймента
	var rsMatch labels.Selector
	for _, rs := range replicaSets.Items {
		if metav1.IsControlledBy(&rs, dep) {
			// Ищем текущую ревизию деплоймента
			if rs.Annotations["deployment.kubernetes.io/revision"] == revision {
				// Сохраняем лейблы этого репликасета
				rsMatch, err = metav1.LabelSelectorAsSelector(rs.Spec.Selector)
				if err != nil {
					return "", err
				}
			}
		}
	}

	// Ищем поды репликасета
	pods, err := k8sRepo.k8sClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: rsMatch.String(),
	})

	sort.Slice(pods.Items, func(i, j int) bool {
		return pods.Items[i].Status.StartTime.Time.Before(pods.Items[j].Status.StartTime.Time)
	})

	if len(pods.Items) == 0 {
		return "", kerrors.NewNotFound(v1.Resource("pods"), "pod")
	}

	return pods.Items[0].Name, nil
}

// CreatePod создает под
func (k8sRepo *k8sRepo) CreatePod(ctx context.Context, name, ns, image string, labels, annotations map[string]string) error {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   ns,
			Annotations: annotations,
			Labels:      labels,
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
	}

	_, err := k8sRepo.k8sClient.CoreV1().Pods(ns).Create(ctx, pod, metav1.CreateOptions{})

	return err
}

// UpdatePod обновляет под
func (k8sRepo *k8sRepo) UpdatePod(ctx context.Context, namespace string, pod *v1.Pod) (*v1.Pod, error) {
	return k8sRepo.k8sClient.CoreV1().Pods(namespace).Update(ctx, pod, metav1.UpdateOptions{})
}

// DeletePod удаляет под
func (k8sRepo *k8sRepo) DeletePod(ctx context.Context, name, namespace string) error {
	return k8sRepo.k8sClient.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (k8sRepo *k8sRepo) CheckPodReadinessProbe(params model.PodProbe) error {
	_, err := (utils.RetryWithResponse{
		N:     10,
		Sleep: 10 * time.Second,
		Fn: func() (string, error) {
			return k8sRepo.PodSendRequest(params.Name, params.Namespace, params.URI, params.Ports)
		},
	}).Start()

	return err
}

// Далем port-forward для нужного пода и отправляем запрос
func (k8sRepo *k8sRepo) PodSendRequest(pod, ns, link, ports string) (result string, err error) {
	log.Printf("Отправить запрос в pod в %s namespace %s, %s\n", pod, ns, link)
	roundTripper, upgrader, err := spdy.RoundTripperFor(k8sRepo.config)
	if err != nil {
		return "", err
	}
	// Запрос на port-forward
	reqURL := k8sRepo.k8sClient.RESTClient().Post().
		Prefix("api/v1").
		Resource("pods").
		Namespace(ns).
		Name(pod).
		SubResource("portforward").URL()

	// Сокет
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodGet, reqURL)

	// Каналы для остановки - stopChan и канал, сигнализирующий о готовности - readyChan
	stopChan, readyChan, errorChan := make(chan struct{}, 1), make(chan struct{}, 1), make(chan error, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	// Подготавливаем объект для port-forward
	forwarder, err := portforward.NewOnAddresses(dialer, []string{"0.0.0.0"}, []string{ports}, stopChan, readyChan, out, errOut)
	if err != nil {
		return "", err
	}

	// Закрывая этот канал, закрывает и лисенеры
	defer close(stopChan)
	// Закрываем port-forward и закрываем лисенеры
	defer forwarder.Close()

	go func() {
		if err = forwarder.ForwardPorts(); err != nil { // Лочится, пока stopChan не будет закрыт, при закрытии stopChan возвращает ошибку
			errorChan <- err
		}
	}()

	// Ждем либо ошибку, либо пока port-forward не станет готовым
	select {
	case err := <-errorChan:
		log.Println("port-forward error", err.Error())
		return "", err
	case <-readyChan: // Этот канал будет закрыт как только port-froward будет готов
		log.Println("port-forward is ready")
	}

	// Проверяем, что нет ошибок, на всякий пожарный
	if len(errOut.String()) != 0 {
		return "", errors.New("error port-forwarding")
	}

	// Отправляем запрос
	if resp, err := http.Get(link); err != nil {
		return "", err
	} else {
		log.Println("STATUS CODE", resp.StatusCode)

		// Проверяем status code
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			log.Println("Все окейси!")
		} else {
			return "", errors.New("response code is not 20x")
		}

		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}
}

// Далем port-forward для нужного пода и отправляем запрос
func (k8sRepo *k8sRepo) PodSendPostRequest(pod, ns, link, ports string, jsonData Data) (result string, err error) {
	log.Printf("Отправить запрос в pod в %s namespace %s, %s\n", pod, ns, link)
	roundTripper, upgrader, err := spdy.RoundTripperFor(k8sRepo.config)
	if err != nil {
		return "", err
	}
	// Запрос на port-forward
	reqURL := k8sRepo.k8sClient.RESTClient().Post().
		Prefix("api/v1").
		Resource("pods").
		Namespace(ns).
		Name(pod).
		SubResource("portforward").URL()

	// Сокет
	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodGet, reqURL)

	// Каналы для остановки - stopChan и канал, сигнализирующий о готовности - readyChan
	stopChan, readyChan, errorChan := make(chan struct{}, 1), make(chan struct{}, 1), make(chan error, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	// Подготавливаем объект для port-forward
	forwarder, err := portforward.NewOnAddresses(dialer, []string{"0.0.0.0"}, []string{ports}, stopChan, readyChan, out, errOut)
	if err != nil {
		return "", err
	}

	// Закрывая этот канал, закрывает и лисенеры
	defer close(stopChan)
	// Закрываем port-forward и закрываем лисенеры
	defer forwarder.Close()

	go func() {
		if err = forwarder.ForwardPorts(); err != nil { // Лочится, пока stopChan не будет закрыт, при закрытии stopChan возвращает ошибку
			errorChan <- err
		}
	}()

	// Ждем либо ошибку, либо пока port-forward не станет готовым
	select {
	case err := <-errorChan:
		log.Println("port-forward error", err.Error())
		return "", err
	case <-readyChan: // Этот канал будет закрыт как только port-froward будет готов
		log.Println("port-forward is ready")
	}

	// Проверяем, что нет ошибок, на всякий пожарный
	if len(errOut.String()) != 0 {
		return "", errors.New("error port-forwarding")
	}

	client := &http.Client{}
	r := bytes.NewReader(jsonData.JsonData)
	// Отправляем запрос
	request, err := http.NewRequest("POST", link, r)
	request.Header.Set("X-PVM-API-KEY", jsonData.ApiKey)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	// Проверяем status code
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		log.Println("Все окейси!")
	} else {
		return "", errors.New("response code is not 20x")
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (k8sRepo *k8sRepo) GetPodLogs(name string, namespace string) (string, error) {
	podLogOpts := v1.PodLogOptions{}
	req := k8sRepo.k8sClient.CoreV1().Pods(namespace).GetLogs(name, &podLogOpts)
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
