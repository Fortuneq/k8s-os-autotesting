package istio

import (
	"context"
	"k8s.io/apimachinery/pkg/types"

	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateSE Создаем ServiceEntry
func (ir *istioRepo) CreateSE(name, ns, host string) error {

	se := &v1beta1.ServiceEntry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: networkingv1beta1.ServiceEntry{
			ExportTo:   []string{"."},
			Hosts:      []string{host},
			Resolution: networkingv1beta1.ServiceEntry_DNS,
			Ports: []*networkingv1beta1.ServicePort{
				{
					Number:   443,
					Name:     "https",
					Protocol: "TLS",
				},
			},
		},
	}

	_, err := ir.istioClient.NetworkingV1beta1().ServiceEntries(ns).Create(context.TODO(), se, metav1.CreateOptions{})

	return err
}

// CreateEmptySpecSE Создаем ServiceEntry с пустым Spec
func (ir *istioRepo) CreateEmptySpecSE(namespace string) error {

	se := &v1beta1.ServiceEntry{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "Empty",
			Namespace: namespace,
		},
		Spec: networkingv1beta1.ServiceEntry{},
	}

	_, err := ir.istioClient.NetworkingV1beta1().ServiceEntries(namespace).Create(context.TODO(), se, metav1.CreateOptions{})

	return err
}
func (ir *istioRepo) CreateSEFromObject(namespace string, obj *v1beta1.ServiceEntry) (*v1beta1.ServiceEntry, error) {
	return ir.istioClient.NetworkingV1beta1().ServiceEntries(namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

// DeleteSE Удаляем ServiceEntry
func (ir *istioRepo) DeleteSE(name, ns string) error {
	return ir.istioClient.NetworkingV1beta1().ServiceEntries(ns).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

func (istioRepo *istioRepo) PatchSERaw(name string, namespace string, patch string, patchOptions metav1.PatchOptions) (*v1beta1.ServiceEntry, error) {
	return istioRepo.istioClient.NetworkingV1beta1().ServiceEntries(namespace).Patch(context.TODO(), name, types.JSONPatchType, []byte(patch), patchOptions, "status")
}
