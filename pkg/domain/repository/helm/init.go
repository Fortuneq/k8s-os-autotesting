package helm

import (
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
)

// Создать клиента для Helm action
func (hr *helmRepo) init(ns string) *action.Configuration {
	if ns == "" {
		ns = "default"
	}
	c := &action.Configuration{}
	// os.Getenv("HELM_DRIVER") - имя Helm draiver, по факту определяет Helm storage (по умолчанию - K8S Secret)
	c.Init(hr.settings.RESTClientGetter(), ns, os.Getenv("HELM_DRIVER"), log.Printf)
	return c
}
