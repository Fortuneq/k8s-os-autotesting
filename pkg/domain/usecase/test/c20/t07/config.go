package t07

import "sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

type T07Config struct {
	AppName    string          `yaml:"appName"`
	Namespace  string          `yaml:"namespace"`
	Workload   params.Workload `yaml:"workload"`
	LowerValue float64         `yaml:"lowerValue"`
	UpperValue float64         `yaml:"upperValue"`
}
