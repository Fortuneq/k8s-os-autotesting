package t02

import "sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

type T02Config struct {
	AppName   string          `yaml:"appName"`
	Namespace string          `yaml:"namespace"`
	Workload  params.Workload `yaml:"workload"`
}
