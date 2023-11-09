package istio

import (
	"sync"

	"istio.io/client-go/pkg/apis/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	istiogo "istio.io/client-go/pkg/clientset/versioned"
	iopv1alpha1add "istio.io/istio/operator/pkg/apis"
	iopv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IstioRepo interface {
	CreateSE(name, ns, host string) error
	CreateEmptySpecSE(namespace string) error
	GetIstioOperator(name, namespace string) (*iopv1alpha1.IstioOperator, error)
	UpdateIstioOperator(iop *iopv1alpha1.IstioOperator) error
	PatchIstioOperator(iop *iopv1alpha1.IstioOperator, patch []byte) error
	GetValidatingWebhookConfiguration(name, namespace string) (*v1.ValidatingWebhookConfiguration, error)
	CreateSEFromObject(namespace string, object *v1beta1.ServiceEntry) (*v1beta1.ServiceEntry, error)
	DeleteSE(name, ns string) error
	CreateIstioOperator(params params.IstioOperatorInstallParams) error
	DeleteIstioOperator(params params.IstioOperatorDeleteParams) error
	PatchSERaw(name string, namespace string, patch string, patchOptions metav1.PatchOptions) (*v1beta1.ServiceEntry, error)
}

var (
	singletone sync.Once
	ir         *istioRepo
)

type istioRepo struct {
	istioClient      *istiogo.Clientset
	k8sRuntimeClient client.WithWatch
}

func CreateNewIstioRepo(config *rest.Config) IstioRepo {
	singletone.Do(func() {
		c1, c2 := createClient(config)
		ir = &istioRepo{
			istioClient:      c1,
			k8sRuntimeClient: c2,
		}
	})

	return ir
}

func createClient(config *rest.Config) (istioClient *istiogo.Clientset, k8sRuntimeClient client.WithWatch) {
	istioClient, err := istiogo.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	k8sRuntimeClient, err = client.NewWithWatch(config, client.Options{})
	if err != nil {
		panic(err.Error())
	}

	if err := iopv1alpha1add.AddToScheme(k8sRuntimeClient.Scheme()); err != nil {
		panic(err.Error())
	}

	return
}
