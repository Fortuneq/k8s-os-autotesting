package helm

import (
	"log"
	"time"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
)

// Устанавливаем релиз
func (hr *helmRepo) Install(params params.HelmInstallParams) error {
	u, err := createInstall(hr.init(params.Namespace), params)

	if err != nil {
		return err
	}

	charts, err := loader.Load(params.Path)
	if err != nil {
		return err
	}

	params.SanitizeValues()
	r, err := u.Run(charts, params.Values)
	if err != nil {
		return err
	}
	log.Println(r.Manifest)

	return nil
}

func createInstall(config *action.Configuration, params params.HelmInstallParams) (*action.Install, error) {
	u := action.NewInstall(config)

	if t, err := time.ParseDuration(params.Timeout); err != nil {
		return nil, err
	} else {
		u.Timeout = t
	}

	u.Atomic = params.Atomic
	u.ReleaseName = params.ReleaseName
	u.Wait = params.Wait
	u.Namespace = params.Namespace

	return u, nil
}
