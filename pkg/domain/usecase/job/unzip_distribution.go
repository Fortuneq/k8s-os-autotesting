package job

import (
	"context"
	"log"

)

func (jobRunner jobRunner) unzipDistribution(context context.Context, jobconfig model.JobConfig) error {
	jobParams, err := prepareConfig[params.UnzipDistributionFilesParams](jobconfig)
	if err != nil {
		return err
	}
	log.Printf("Unzipping %s to %s", jobParams.DistributionSavePath, jobParams.UnzippedDistributionPath)
	if err := utils.Unzip(jobParams.DistributionSavePath, jobParams.UnzippedDistributionPath); err != nil {
		return err
	}
	log.Println("Unzipped")
	return err
}
