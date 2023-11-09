package job

import (
	"context"

	"gopkg.in/yaml.v2"

)

func (jobRunner *jobRunner) deleteTestApp(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		paramsFromConfig, err := config.ReadConfigFromByteArray[params.TestAppParams](b)
		if err != nil {
			return err
		}

		if paramsFromConfig.Workload == params.DeploymentWorkload {
			return jobRunner.k8sRepo.DeleteTestAppDeployment(ctx, paramsFromConfig)
		} else {
			return jobRunner.k8sRepo.DeleteTestAppPod(ctx, paramsFromConfig)
		}
	} else {
		return err
	}
}
