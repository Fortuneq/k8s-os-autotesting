package t09


type T09Config struct {
	AppName       string          `yaml:"appName"`
	LoadNamespace string          `yaml:"loadNamespace"`
	Workload      params.Workload `yaml:"workload"`
	Memload       int             `yaml:"memload"`
	Timeout       int             `yaml:"timeout"`
}
