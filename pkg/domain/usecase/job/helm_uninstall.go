package job

import (
	"context"

	"gopkg.in/yaml.v2"

)

func (jobRunner *jobRunner) helmUninstall(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[params.UninnstallParams](b)
		if err != nil {
			return err
		}

		return jobRunner.helmRepo.Uninstall(params)
	} else {
		return err
	}
}
