package helm

import (
	"log"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"
)

// Проверяем есть ли указанный релиз в неймспейсе
func (hr *helmRepo) Get(params params.HelmInstallParams) (bool, error) {
	// Создаем клиент
	get := action.NewGet(hr.init(params.Namespace))
	// Получаем релиз
	if _, err := get.Run(params.ReleaseName); err != nil && (err == driver.ErrReleaseNotFound || err == driver.ErrNoDeployedReleases) {
		// Не нашли
		log.Printf("Release %s is not found, new release will be created\n", params.ReleaseName)
		return false, nil
	} else if err != nil {
		// Ошибка (отличная от отсуствия релиза)
		return false, err
	}

	// Попался пездюк!
	return true, nil
}
