package job

import (
	"context"
	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/yaml"
	"strings"

	"gopkg.in/yaml.v3"


)

type JobRunner interface {
	// todo: копипаста много
	RunJob(jobConfig model.JobConfig) error
}

type jobRunner struct { //мб ренейм
	jobsMap   map[string]model.JobRunner
	helmRepo  helm.HelmRepo
	istioRepo istio.IstioRepo
	k8sRepo   kubernetes.K8SRepo
}

func NewJobRunner() JobRunner {
	config := utils.KubeRestConfig()
	jobRunner := &jobRunner{
		helmRepo:  helm.CreateNewHelmRepo(),
		istioRepo: istio.CreateNewIstioRepo(config),
		k8sRepo:   kubernetes.CreateNewK8SRepo(config),
	}
	jobRunner.jobsMap = map[string]model.JobRunner{
		"helm-install":                         jobRunner.helmInstall,
		"helm-uninstall":                       jobRunner.helmUninstall,
		"helm-upgrade":                         jobRunner.helmUpgrade,
		"create-namespace":                     jobRunner.createNamespace, // todo ?
		"delete-namespace":                     jobRunner.deleteNamespace,
		"create-pull-secret":                   jobRunner.createPullSecret,
		"create-operator":                      jobRunner.createIstioOperator,
		"download-distribution":                jobRunner.downloadDistribution,
		"unzip-distribution":                   jobRunner.unzipDistribution,
		"clear-distribution-files":             jobRunner.creatDistributionFiles,
		"create-test-app":                      jobRunner.createTestApp,
		"delete-test-app":                      jobRunner.deleteTestApp,
		"create-network-attachemtn-definition": jobRunner.createNAD,
		"delete-operator":                      jobRunner.deleteIOP,
	}
	return jobRunner
}

func (jobRunner *jobRunner) RunJob(jobConfig model.JobConfig) error {
	if job := jobRunner.jobsMap[jobConfig.Name]; job != nil {
		//log.Printf("Config: %+v\n", jobConfig)
		if err := job(context.Background(), jobConfig); err != nil {
			return err
		}
	} else {
		sb := strings.Builder{}
		for k := range jobRunner.jobsMap {
			sb.WriteString(k)
			sb.WriteRune('\n')
		}
		return errors.New(fmt.Sprintf("No such job %s. Available jobs: %v", jobConfig.Name, sb.String()))
	}

	return nil
}

func prepareConfig[T any](jobConfig model.JobConfig) (T, error) {
	if b, err := yaml.Marshal(jobConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T](b)
		if err = yaml.Unmarshal(b, &params); err != nil {
			return params, err
		}
		return params, nil
	} else {
		var nilParams T
		return nilParams, err
	}
}
