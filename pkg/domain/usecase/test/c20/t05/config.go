package t05

import "sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

type T05Config struct {
	AppName       string          `yaml:"appName"`
	LoadNamespace string          `yaml:"loadNamespace"`
	TestNamespace string          `yaml:"testNamespace"`
	Workload      params.Workload `yaml:"workload"`
	Cpuload       int             `yaml:"cpuLoad"`
	Timeout       int             `yaml:"timeout"`
}
