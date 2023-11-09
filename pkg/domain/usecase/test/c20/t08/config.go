package t08



type T08Config struct {
	AppName   string          `yaml:"appName"`
	Namespace string          `yaml:"namespace"`
	Workload  params.Workload `yaml:"workload"`
	PodCount  int             `yaml:"podCount"`
	Timeout   int             `yaml:"timeout"`
}
