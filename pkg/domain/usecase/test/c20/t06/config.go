package t06

import "sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"

type T06Config struct {
	AppName   string          `yaml:"appName"`
	Namespace string          `yaml:"namespace"`
	Workload  params.Workload `yaml:"workload"`
	ApiKey    string          `yaml:"apiKey"`
	SqlQuery  string          `yaml:"sqlQuery"`
}
