// todo:

package t01

import (
	"context"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type T01Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T01Config
}

func (ises *T01Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T01Suite) New() model.RunnableTest {
	return &T01Suite{}
}

func (ises *T01Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T01Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T01Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T01Config](b)
		if err != nil {
			t.Fatal(err.Error())
		}
		ises.testConfig = params

		if err = yaml.Unmarshal(b, &params); err != nil {
			t.Fatal(err.Error())
		}
	} else {
		t.Fatal(err.Error())
	}
}

func (ises *T01Suite) Test(t provider.T) {
	t.Title("Тестирование Готовности оператора")
	t.Description("Проверяем, то , что поднялся и готов autoscaler operator")
	ises.prepareConfig(t)

	ises.checkTestAppReadinessProbe(t)
}

func (ises *T01Suite) checkTestAppReadinessProbe(t provider.T) {
	t.WithNewStep("Проверяем Readiness probe proxy autoscaler operator", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем pod operator по названию деплоймента")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.AppName, ises.testConfig.Namespace)
			if err != nil {
				t.Fatal(err.Error())
			}

			readinessProbeStep := allure.NewSimpleStep("Проверяем Readiness probe")
			readinessProbeStep.WithParent(ctx.CurrentStep())
			if err = ises.k8sRepo.
				CheckPodReadinessProbe(model.PodProbe{
					Name:      podName,
					Namespace: ises.testConfig.Namespace,
					URI:       "http://0.0.0.0:8082/readyz",
					Ports:     "8082:8082",
				}); err != nil {
				t.Fatal(err.Error())
			}
			logStr := fmt.Sprintf("pod %s readiness probe прошла успешно, status code = 200", podName)
			readinessProbeStep.WithAttachments(allure.NewAttachment("Статус readiness probe", allure.Text, []byte(logStr)))
		} else {
			readinessProbeStep := allure.NewSimpleStep("Проверяем Readiness probe")
			readinessProbeStep.WithParent(ctx.CurrentStep())
			if err := ises.k8sRepo.CheckPodReadinessProbe(model.PodProbe{
				Name:      ises.testConfig.AppName,
				Namespace: ises.testConfig.Namespace,
				URI:       "http://0.0.0.0:8082/readyz",
				Ports:     "8082:8082",
			}); err != nil {
				t.Fatal(err.Error())
			}
			logStr := fmt.Sprintf("pod %s readiness probe прошла успешно, status code = 200", ises.testConfig.AppName)
			readinessProbeStep.WithAttachments(allure.NewAttachment("Статус readiness probe", allure.Text, []byte(logStr)))
		}
	})
}
