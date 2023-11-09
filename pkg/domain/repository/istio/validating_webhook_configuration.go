package istio

import (
	"context"

	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (istioRepo *istioRepo) GetValidatingWebhookConfiguration(name, namespace string) (*v1.ValidatingWebhookConfiguration, error) {
	var vwc v1.ValidatingWebhookConfiguration

	if err := istioRepo.k8sRuntimeClient.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &vwc, &client.GetOptions{}); err != nil {
		return &v1.ValidatingWebhookConfiguration{}, err
	}
	return &vwc, nil
}
