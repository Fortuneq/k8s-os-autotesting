package helm

import (
	"helm.sh/helm/v3/pkg/cli"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"
	"sbet-tech.com/synapse/istio-se/allure/pkg/utils"
)

type HelmRepo interface {
	Install(params.HelmInstallParams) error
	Uninstall(uninnstallParams params.UninnstallParams) error
	InstallOrUpgrade(params params.HelmInstallParams) error
	Upgrade(params params.HelmInstallParams) error
	Get(params params.HelmInstallParams) (bool, error)
}

type helmRepo struct {
	settings *cli.EnvSettings
}

func CreateNewHelmRepo() HelmRepo {
	return &helmRepo{settings: utils.HelmClient()}
}
