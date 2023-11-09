package params

type Workload string

const (
	DeploymentWorkload Workload = "deployment"
	PodWorkload        Workload = "pod"
)

var NewIDValue int64 = 1234

type TestAppParams struct {
	Namespace   string            `yaml:"namespace"`
	Image       string            `yaml:"image"`
	AppName     string            `yaml:"appName"`
	Workload    Workload          `yaml:"workload"`
	Labels      map[string]string `yaml:"labels"`
	Annotations map[string]string `yaml:"annotations"`
}
