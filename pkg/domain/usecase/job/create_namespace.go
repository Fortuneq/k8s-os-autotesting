package job

import (
	"context"
)

func (jobRunner *jobRunner) createNamespace(ctx context.Context, jobconfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobconfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[params.NamespaceParams](b)
		if err != nil {
			return err
		}

		return jobRunner.k8sRepo.CreateNamespace(ctx, params)
	} else {
		return err
	}
}
