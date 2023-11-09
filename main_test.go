package main

import (
	"log"
	"testing"
)

// todo run golangci-lint
func TestRunner(t *testing.T) {
	utils.InitHelmClient(utils.GetKubeconfigPath())

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Println("Config Path", confPath)
	log.Println("Kubeconf Path", utils.GetKubeconfigPath())

	tasks, err := config.ReadConfigFromFile[model.Task](confPath)
	if err == nil {
		tr := taskrunner.NewTaskRunner(tasks, t)
		if err := tr.RunTask(); err != nil {
			t.Fatalf("error running: %v", err.Error())
		}
	} else {
		t.Fatalf("error reading config: %v", err.Error())
	}
}
