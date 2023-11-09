package config

import (
	"log"
	"os"
)

func ReadConfigFromFile[T any](taskConfigFilePath string) (T, error) {
	yamlFile, err := os.ReadFile(taskConfigFilePath)
	if err != nil {
		var nilObject T
		return nilObject, err
	}
	return ReadConfigFromByteArray[T](yamlFile)
}

func ReadConfigFromString[T any](config string) (T, error) {
	return ReadConfigFromByteArray[T]([]byte(config))
}

func ReadConfigFromByteArray[T any](config []byte) (T, error) {
	var conf T
	if err := yaml.Unmarshal(config, &conf); err == nil {
		//log.Printf("Unmarshal config %#v", conf)
		return conf, err
	} else {
		log.Println("Error reading config")
		return conf, err
	}
}
