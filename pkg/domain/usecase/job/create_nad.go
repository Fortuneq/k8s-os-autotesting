package job

import (
	"context"

	"gopkg.in/yaml.v3"
)

func (jobRunner *jobRunner) createNAD(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		paramsFromConfig, err := config.ReadConfigFromByteArray[params.NADParams](b)
		if err != nil {
			return err
		}

		return jobRunner.k8sRepo.CreateNAD(paramsFromConfig)
	} else {
		return err
	}
}
