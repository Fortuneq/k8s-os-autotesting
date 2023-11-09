package params

type IstioOperatorDeleteParams struct {
	// Имя конфига
	Name string `yaml:"name"`
	// Проект конфига
	Namespace string `yaml:"namespace"`
}
