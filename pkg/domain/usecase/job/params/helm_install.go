package params

import (
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/kustomize/kyaml/yaml"
)

// Параметры обновления/установки для Helm релиза
type HelmInstallParams struct {
	// Проект Helm, используется для Helm Storage и в качестве Release.Namespace
	Namespace string `yaml:"namespace"`
	// Имя Helm релиза
	ReleaseName string `yaml:"releaseName"`
	// В значении true Helm будет дожидаться установки ресурсов, актуально для Workload (время по умолчанию 5 минут)
	Wait bool `yaml:"wait"`
	// В значении true Helm установит либо все конфиги, либо в случае возникнование ошибки установки хотя бы одного конфига, откатит релиз к прошлой версии
	Atomic bool `yaml:"atomic"`
	// Время ожидания установки ресурсов, работает в паре в параметром Wait
	Timeout string `yaml:"timeout"`
	// Путь до Helm charts
	Path string `yaml:"path"`
	// Актуально Upgrade. В значении true, если релиз не был установлен, то Upgrade его установит
	Install bool `yaml:"install"`
	// Переопределение Values из Helm Charts
	Values map[string]interface{} `yaml:"values"`
}

// SanitizeValues - костыль, чтобы не ловить от хельма тайп мисматч
func (p *HelmInstallParams) SanitizeValues() {
	if p.Values != nil {
		if b, err := yaml.Marshal(p.Values); err != nil {
			panic(err.Error())
		} else {
			if v, err := chartutil.ReadValues(b); err != nil {
				panic(err.Error())
			} else {
				p.Values = v
			}
		}
	}
}
