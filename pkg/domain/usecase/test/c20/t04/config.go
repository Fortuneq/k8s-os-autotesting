package t04

import "sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

type T04Config struct {
	LoadAppName   string          `yaml:"loadAppName"`
	TestAppName   string          `yaml:"testAppName"`
	LoadNamespace string          `yaml:"loadNamespace"`
	TestNamespace string          `yaml:"testNamespace"`
	Workload      params.Workload `yaml:"workload"`
	PodCount      int             `yaml:"podCount"`
}
