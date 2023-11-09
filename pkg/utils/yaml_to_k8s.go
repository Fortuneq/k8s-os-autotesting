package utils

import (
	"bytes"

	k8Yaml "k8s.io/apimachinery/pkg/util/yaml"
)

func ReadYamlToObject(file []byte, obj any) error {
	dec := k8Yaml.NewYAMLOrJSONDecoder(bytes.NewReader(file), 1000)
	return dec.Decode(obj)
}
