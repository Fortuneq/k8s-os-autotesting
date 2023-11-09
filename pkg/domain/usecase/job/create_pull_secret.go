package job

import (
	"context"
)

func (jobRunner *jobRunner) createPullSecret(ctx context.Context, jobconfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobconfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[params.PullSecretParams](b)
		if err != nil {
			return err
		}

		return jobRunner.k8sRepo.CreatePullSecret(ctx, params)
	} else {
		return err
	}
}
