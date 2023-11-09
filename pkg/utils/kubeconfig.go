package utils

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/spf13/pflag"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	singletoneKC   sync.Once
	config         *rest.Config
	kubeconfigPath string
)

func AddKubeConfigFlag(fs *pflag.FlagSet) {
	if home := homedir.HomeDir(); home != "" {
		log.Println("Kubeconfig in home")
		fs.StringVar(&kubeconfigPath, "kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		log.Println("Kubeconfig heeds flag")
		fs.StringVar(&kubeconfigPath, "kubeconfig", "", "absolute path to the kubeconfig file")
	}
}

func GetKubeconfigPath() string {
	return kubeconfigPath
}

func KubeRestConfig() *rest.Config {
	singletoneKC.Do(initConfig)
	return config
}

func initConfig() {
	// Мы в кластере
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		if err != rest.ErrNotInCluster {
			panic(err.Error())
		} else {
			log.Printf("Not in cluster")

			// Используем текущий контекст
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}
