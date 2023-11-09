package params

// Параметры обновления/установки для Helm релиза
type UninnstallParams struct {
	// Проект Helm, используется для Helm Storage и в качестве Release.Namespace
	Namespace string `yaml:"namespace"`
	// Имя Helm релиза
	ReleaseName string `yaml:"releaseName"`
	// В значении true Helm будет дожидаться установки ресурсов, актуально для Workload (время по умолчанию 5 минут)
	Wait bool `yaml:"wait"`

	Timeout string `yaml:"timeout"`

	// В значении true Helm установит либо все конфиги, либо в случае возникнование ошибки установки хотя бы одного конфига, откатит релиз к прошлой версии
}
