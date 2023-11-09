package job

import (
	"context"

	"gopkg.in/yaml.v2"
)

func (jobRunner *jobRunner) helmUpgrade(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[params.HelmInstallParams](b)
		if err != nil {
			return err
		}

		if params.Install {
			return jobRunner.helmRepo.InstallOrUpgrade(params)
		} else {
			return jobRunner.helmRepo.Upgrade(params)
		}
	} else {
		return err
	}
}
