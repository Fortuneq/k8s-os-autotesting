package helm

import (
	"log"
	"time"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// Обновляем релиз
func (hr *helmRepo) Upgrade(params params.HelmInstallParams) error {
	u, err := createUpgrade(hr.init(params.Namespace), params)
	if err != nil {
		return err
	}

	charts, err := loader.Load(params.Path)
	if err != nil {
		return err
	}

	params.SanitizeValues()
	r, err := u.Run(params.ReleaseName, charts, params.Values)
	if err != nil {
		return err
	}
	log.Println(r.Manifest)

	return nil
}

func createUpgrade(config *action.Configuration, params params.HelmInstallParams) (*action.Upgrade, error) {
	u := action.NewUpgrade(config)

	if t, err := time.ParseDuration(params.Timeout); err != nil {
		return nil, err
	} else {
		u.Timeout = t
	}

	u.Atomic = params.Atomic
	u.Wait = params.Wait
	u.Namespace = params.Namespace

	return u, nil
}
