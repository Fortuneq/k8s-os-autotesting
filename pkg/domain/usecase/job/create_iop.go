package job

import (
	"context"

	"gopkg.in/yaml.v3"

)

func (jobRunner *jobRunner) createIstioOperator(ctx context.Context, jobconfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobconfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[params.IstioOperatorInstallParams](b)
		if err != nil {
			return err
		}
		return jobRunner.istioRepo.CreateIstioOperator(params)
	} else {
		return err
	}
}
