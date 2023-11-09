package helm

import (
	"time"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	"helm.sh/helm/v3/pkg/action"
)

// Устанавливаем релиз
func (hr *helmRepo) Uninstall(params params.UninnstallParams) error {
	u, err := createUninstall(hr.init(params.Namespace), params)

	if err != nil {
		return err
	}

	_, err = u.Run(params.ReleaseName)
	if err != nil {
		return err
	}

	return nil
}

func createUninstall(config *action.Configuration, params params.UninnstallParams) (*action.Uninstall, error) {
	u := action.NewUninstall(config)

	if t, err := time.ParseDuration(params.Timeout); err != nil {
		return nil, err
	} else {
		u.Timeout = t
	}

	u.Wait = params.Wait

	return u, nil
}
