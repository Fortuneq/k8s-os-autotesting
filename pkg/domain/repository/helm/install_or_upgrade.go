package helm

import (
	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"
)

// Установить или обновить Helm релиз
func (hr *helmRepo) InstallOrUpgrade(params params.HelmInstallParams) error {

	// Проверяем есть ли релиз
	if found, err := hr.Get(params); err != nil {
		return err
	} else if !found {
		// Если нет, то ставим
		return hr.Install(params)
	} else {
		// Если есть, то обновляем
		return hr.Upgrade(params)
	}
}
