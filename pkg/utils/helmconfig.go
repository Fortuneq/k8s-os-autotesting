package utils

import (
	"sync"

	"helm.sh/helm/v3/pkg/cli"
)

var (
	singletoneHC sync.Once
	helmEnv      *cli.EnvSettings
)

func HelmClient() *cli.EnvSettings {
	if helmEnv == nil {
		panic("Helm client is not initialized. You shoul call method InitHelmClient first")
	}
	return helmEnv
}

func InitHelmClient(kubeconfigPath string) {
	singletoneHC.Do(func() {
		helmEnv = cli.New()
		helmEnv.KubeConfig = kubeconfigPath
	})
}
