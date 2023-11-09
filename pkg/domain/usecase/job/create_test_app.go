package job

import (
	"context"

	"gopkg.in/yaml.v2"

)

func (jobRunner *jobRunner) createTestApp(ctx context.Context, jobConfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		paramsFromConfig, err := config.ReadConfigFromByteArray[params.TestAppParams](b)
		if err != nil {
			return err
		}

		if paramsFromConfig.Workload == params.DeploymentWorkload {
			return jobRunner.k8sRepo.CreateTestAppDeployment(ctx, paramsFromConfig)
		} else {
			return jobRunner.k8sRepo.CreateTestAppPod(ctx, paramsFromConfig)
		}
	} else {
		return err
	}
}
