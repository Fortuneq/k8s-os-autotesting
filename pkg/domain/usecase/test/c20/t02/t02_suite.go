// todo:

package t02

import (
	"context"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type T02Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T02Config
}

func (ises *T02Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T02Suite) New() model.RunnableTest {
	return &T02Suite{}
}

func (ises *T02Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T02Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T02Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T02Config](b)
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

func (ises *T02Suite) Test(t provider.T) {
	t.Title("Тестирование Readyness probe тестового приложения")
	t.Description("Проверяем, что под тест приложения поднялся")
	ises.prepareConfig(t)

	ises.checkTestAppReady(t)
}

func (ises *T02Suite) checkTestAppReady(t provider.T) {
	t.WithNewStep("Проверяем Readiness тестового приложения", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем pod тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.AppName, ises.testConfig.Namespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			logStr := fmt.Sprintf("pod %s найден, status code = 200", podName)
			getPodStep.WithAttachments(allure.NewAttachment("Нашли под тест приложения", allure.Text, []byte(logStr)))
		}
	})
}
