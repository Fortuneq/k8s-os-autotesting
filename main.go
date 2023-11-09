package main

/**
* Суть в том, что он не хотел хавать флаги в тестах, пока они не будут объявлены в исполняемом коде
**/

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
)

var (
	confPath string
)

func init() {
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	utils.AddKubeConfigFlag(flag.CommandLine)

	flag.StringVar(&confPath, "config-path", "./test.yaml", "Path to config file")

	flag.Parse()
}

func main() {

}
