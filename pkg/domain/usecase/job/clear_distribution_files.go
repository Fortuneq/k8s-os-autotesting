package job

import (
	"context"
	"log"
	"os"

)

func (jobRunner *jobRunner) creatDistributionFiles(ctx context.Context, jobconfig model.JobConfig) error {
	jobParams, err := prepareConfig[params.ClearDistributionFilesParams](jobconfig)
	if err != nil {
		return err
	}
	log.Println("Removing " + jobParams.DistributionSavePath)
	if err := os.RemoveAll(jobParams.DistributionSavePath); err != nil {
		return err
	}
	log.Println("Removing " + jobParams.UnzippedDistributionPath)
	if err := os.RemoveAll(jobParams.UnzippedDistributionPath); err != nil {
		return err
	}
	log.Println("Removed")
	return err
}
