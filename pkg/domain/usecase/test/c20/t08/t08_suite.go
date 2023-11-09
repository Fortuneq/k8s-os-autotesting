package t08

import (
	"context"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"k8s.io/apimachinery/pkg/util/yaml"
	"time"
)

type T08Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T08Config
}

func (ises *T08Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T08Suite) New() model.RunnableTest {
	return &T08Suite{}
}

func (ises *T08Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T08Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T08Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T08Config](b)
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

func (ises *T08Suite) Test(t provider.T) {
	t.Title("Тестирование Равенства кол-ва подов к ожидаемому")
	t.Description("Проверяем, что количество подов такое, какое мы ожидаем")
	ises.prepareConfig(t)

	ises.checkTestAppPodCount(t)
}

func (ises *T08Suite) checkTestAppPodCount(t provider.T) {
	t.WithNewStep("Проверяем что после скейлинга/даунскейлинга осталось конкретное количество подов", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podNames, err := ises.k8sRepo.GetAllDeploymentPodNames(context.TODO(), ises.testConfig.AppName, ises.testConfig.Namespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			for i := 0; i < 5; i++ {
				if len(podNames) > ises.testConfig.PodCount {
					time.Sleep(30 * time.Second)
				}
				break
			}
			if len(podNames) > ises.testConfig.PodCount {
				t.Fatal(fmt.Errorf("количество реплик %v > желаемое %v", len(podNames), ises.testConfig.PodCount))
			}
			logStr := fmt.Sprintf(" %v podов осталось, status code = 200", len(podNames))
			getPodStep.WithAttachments(allure.NewAttachment("Количество подов равно ожидаемому", allure.Text, []byte(logStr)))
		}
	})
}
