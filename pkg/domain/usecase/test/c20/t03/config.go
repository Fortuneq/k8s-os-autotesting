package t03

import "sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

type T03Config struct {
	AppName       string          `yaml:"appName"`
	Namespace     string          `yaml:"namespace"`
	TestNamespace string          `yaml:"testNamespace"`
	Workload      params.Workload `yaml:"workload"`
	Cpuload       int             `yaml:"cpuLoad"`
}
