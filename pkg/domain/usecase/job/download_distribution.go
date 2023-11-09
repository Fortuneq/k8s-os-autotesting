package job

import (
	"context"
	"log"

)

func (jobRunner *jobRunner) downloadDistribution(context context.Context, jobconfig model.JobConfig) error {
	jobParams, err := prepareConfig[params.DownloadDistributionParams](jobconfig)
	if err != nil {
		return err
	}
	log.Println("Downloading from: " + jobParams.DistributionURL)
	path, err := utils.LoadFileFromURLToPath(jobParams.DistributionURL, jobParams.DistributionSavePath, jobParams.Username, jobParams.Password)
	if err != nil {
		return err
	}
	log.Println("Saved to: " + path)
	return err
}
