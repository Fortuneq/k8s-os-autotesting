package kubernetes

import (
	"context"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/scale"
	"sync"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/model"
	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8SRepo interface {
	CreateNamespace(cxt context.Context, params params.NamespaceParams) error
	DeleteNamespace(cxt context.Context, params params.NamespaceParams) error
	CheckNamespace(ns string) error
	PodSendRequest(pod, ns, link, ports string) (string, error)
	PodSendPostRequest(pod, ns, link, ports string, jsonData Data) (string, error)
	DeployDeployment(ctx context.Context, appName, namespace, image string, labels, annotations map[string]string) error
	CreateDeployment(ctx context.Context, name, ns, image string, labels, annotations map[string]string) error
	UpdateDeployment(ctx context.Context, namespace string, deployment *appsV1.Deployment) (*appsV1.Deployment, error)
	DeleteDeployment(ctx context.Context, name, ns string) error
	DeployPod(ctx context.Context, name, namespace, image string, labels, annotations map[string]string) error
	CreatePod(ctx context.Context, name, ns, image string, labels, annotations map[string]string) error
	GetPod(ctx context.Context, name, namespace string) (*coreV1.Pod, error)
	GetPods(ctx context.Context, namespace string) (*coreV1.PodList, error)
	CheckPodReadinessProbe(params model.PodProbe) error // todo: по факту выполняет функцию PodSendRequest, нужно будет сюда вбить надо проб
	UpdatePod(ctx context.Context, namespace string, pod *coreV1.Pod) (*coreV1.Pod, error)
	DeletePod(ctx context.Context, name, ns string) error
	GetAllDeploymentPodNames(ctx context.Context, deployment, namespace string) ([]string, error)
	GetCurrentUsageMetrics(selector string, namespace string) (result *CpuMemoryPodsCount)
	GetLimits(selector string, namespace string) (result *CpuMemory)
	GetPodName(ctx context.Context, deployment, namespace string) (string, error)
	GetClusterRole(name string) (*rbacV1.ClusterRole, error)
	GetClusterRoleBinding(name string) (*rbacV1.ClusterRoleBinding, error)
	GetDeployment(ctx context.Context, name, namespace string) (*appsV1.Deployment, error)
	GetDeploymentName(namespace string, selector string) (string, error)
	GetRole(name, namespace string) (*rbacV1.Role, error)
	GetService(name, namespace string) (*coreV1.Service, error)
	GetServiceAccount(name, namespace string) (*coreV1.ServiceAccount, error)
	GetDaemonset(name, ns string) (*appsV1.DaemonSet, error)
	CreatePullSecret(ctx context.Context, params params.PullSecretParams) error
	CreateTestAppDeployment(ctx context.Context, params params.TestAppParams) error
	DeleteTestAppDeployment(ctx context.Context, params params.TestAppParams) error
	CreateTestAppPod(ctx context.Context, params params.TestAppParams) error
	DeleteTestAppPod(ctx context.Context, params params.TestAppParams) error
	GetPodLogs(name string, namespace string) (string, error)
	CreateNAD(params params.NADParams) error
}

var (
	singletone sync.Once
	kr         *k8sRepo
)

type k8sRepo struct {
	k8sClient        *kubernetes.Clientset
	config           *rest.Config
	k8sRuntimeClient client.Client
	metricsClient    *metrics.Clientset
	scalerGetter     scale.ScalesGetter
}

func CreateNewK8SRepo(config *rest.Config) K8SRepo {
	singletone.Do(func() {
		c1, c2, c3 := createClient(config)
		kr = &k8sRepo{
			k8sClient:        c1,
			config:           config,
			k8sRuntimeClient: c2,
			metricsClient:    c3,
		}
	})

	return kr
}

func createClient(config *rest.Config) (k8sClient *kubernetes.Clientset, k8sRuntimeClient client.Client, metricsClient *metrics.Clientset) {
	// Создаем clientset
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	k8sRuntimeClient, err = client.New(config, client.Options{})
	if err != nil {
		panic(err.Error())
	}

	metricsClient, err = metrics.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	if err := netv1.AddToScheme(k8sRuntimeClient.Scheme()); err != nil {
		panic(err.Error())
	}

	return
}
