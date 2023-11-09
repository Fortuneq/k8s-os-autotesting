package job

import (
	"context"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func (jobRunner *jobRunner) deleteNamespace(ctx context.Context, jobconfig model.JobConfig) error {
	if b, err := yaml.Marshal(jobconfig.Params); err == nil {
		var params params.NamespaceParams
		if err = yaml.Unmarshal(b, &params); err != nil {
			return err
		}

		return jobRunner.k8sRepo.DeleteNamespace(ctx, params)
	} else {
		return err
	}
}
