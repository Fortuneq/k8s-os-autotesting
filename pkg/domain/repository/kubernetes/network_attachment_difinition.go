package kubernetes

import (
	"context"
	"log"

	netv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (k8sRepo *k8sRepo) CreateNAD(params params.NADParams) error {
	var nad netv1.NetworkAttachmentDefinition

	if err := k8sRepo.k8sRuntimeClient.Get(context.TODO(), types.NamespacedName{
		Namespace: params.Namespace,
		Name:      params.Name,
	}, &nad, &client.GetOptions{}); err != nil {
		newNad := &netv1.NetworkAttachmentDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      params.Name,
				Namespace: params.Namespace,
			},
		}

		if kerrors.IsNotFound(err) {
			return k8sRepo.k8sRuntimeClient.Create(context.TODO(), newNad, &client.CreateOptions{})
		} else {
			return err
		}
	} else {
		log.Printf("NetworkAttachmentDefinition %s уже есть в проекте %s", params.Name, params.Namespace)
	}

	return nil
}
