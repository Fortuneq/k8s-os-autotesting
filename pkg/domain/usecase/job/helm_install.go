package job

import (
	"context"

	"gopkg.in/yaml.v2"

)

func (jobRunner *jobRunner) helmInstall(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[params.HelmInstallParams](b)
		if err != nil {
			return err
		}

		return jobRunner.helmRepo.Install(params)
	} else {
		return err
	}
}
