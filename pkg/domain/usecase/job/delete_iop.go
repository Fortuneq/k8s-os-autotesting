package job

import (
	"context"

	"gopkg.in/yaml.v3"

)

func (jobRunner *jobRunner) deleteIOP(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		paramsFromConfig, err := config.ReadConfigFromByteArray[params.IstioOperatorDeleteParams](b)
		if err != nil {
			return err
		}

		return jobRunner.istioRepo.DeleteIstioOperator(paramsFromConfig)
	} else {
		return err
	}
}
