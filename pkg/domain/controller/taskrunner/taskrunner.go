package taskrunner

import (
	"context"
	"log"
	"testing"

	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/model"
	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job"
	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/test"
)

type TaskRunner interface {
	RunTask() error
}

type runner struct {
	ctx           context.Context
	task          model.Task
	k8STestRunner test.K8sTestRunner
	jobRunner     job.JobRunner
}

func NewTaskRunner(task model.Task, provider *testing.T) TaskRunner {
	return &runner{
		k8STestRunner: test.NewK8sTestRunner(provider),
		jobRunner:     job.NewJobRunner(),
		task:          task,
	}
}

func (runner *runner) RunTask() error {
	if err := runner.runJobs(runner.task.Jobs); err != nil {
		return err
	}
	if err := runner.runTests(runner.task.Tests); err != nil {
		return err
	}
	if err := runner.runJobs(runner.task.Finalizers); err != nil {
		return err
	}
	return nil
}

func (runner *runner) runJobs(jobs []model.JobConfig) error {
	log.Println("###############")
	log.Println("Running jobs")
	var err error
	for _, job := range jobs {
		log.Println("Run job", job.Name)
		if job.Params != nil {
			// todo: mask
			//log.Println("Params", job.Params.(map[any]any))
		}

		err = runner.jobRunner.RunJob(job)

		if err != nil {
			log.Println("Error on executing job", job.Name, err.Error())
			return err //todo
		} else {
			log.Println("Job is complete")
		}
	}
	log.Println("Running jobs finished")
	log.Println("###############")
	return nil
}

func (runner *runner) runTests(tests []model.TestConfig) error {
	log.Println("###############")
	log.Println("Running tests")
	for _, testConfig := range tests {
		log.Println("Run test", testConfig.Name)
		if err := runner.k8STestRunner.RunTest(testConfig); err != nil {
			log.Println("Error on executing job", testConfig.Name, err.Error())
			return err
		}
	}
	log.Println("Running tests finished")
	log.Println("###############")
	return nil
}
