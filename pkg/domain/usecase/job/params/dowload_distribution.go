package params

type DownloadDistributionParams struct {
	DistributionURL      string `yaml:"distributionURL"`
	DistributionSavePath string `yaml:"distributionSavePath"`
	Username             string `yaml:"username"`
	Password             string `yaml:"password"`
}
