package t07

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
	"time"
)

type T07Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T07Config
}

func (ises *T07Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T07Suite) New() model.RunnableTest {
	return &T07Suite{}
}

func (ises *T07Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T07Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T07Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T07Config](b)
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

func (ises *T07Suite) Test(t provider.T) {
	t.Title("Меняем значения которые отдает мок монитора по query запросу")
	ises.prepareConfig(t)

	ises.requestTestAppMonitor(t)
}

func (ises *T07Suite) requestTestAppMonitor(t provider.T) {
	t.WithNewStep("Отправляем запрос к тестовому приложению моку монитора", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.AppName, ises.testConfig.Namespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			changeLoadStep := allure.NewSimpleStep("Запрос по моку монитора для смены рандомайз значений")
			values := map[string]float64{"lowerValue": ises.testConfig.LowerValue, "upperValue": ises.testConfig.UpperValue}
			jsonValue, _ := json.Marshal(values)
			string, err := ises.k8sRepo.PodSendPostRequest(podName, ises.testConfig.Namespace,
				"http://localhost:5000/set_value", "5000:3000",
				kubernetes.Data{JsonData: jsonValue})
			if err != nil {
				fmt.Println(string)
				t.Fatal(err.Error())
			}
			changeLoadStep.WithParent(getPodStep)
			time.Sleep(75 * time.Second)
			logStr := fmt.Sprintf("Запрос прошел")
			getPodStep.WithAttachments(allure.NewAttachment("Все запросы к поду прошли", allure.Text, []byte(logStr)))
		}
	})
}
