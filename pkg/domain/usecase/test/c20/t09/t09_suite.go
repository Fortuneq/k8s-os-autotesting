package t09

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

type T09Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T09Config
}

func (ises *T09Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T09Suite) New() model.RunnableTest {
	return &T09Suite{}
}

func (ises *T09Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T09Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T09Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T09Config](b)
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

func (ises *T09Suite) Test(t provider.T) {
	t.Title("Изменение нагрузки в mock приложении для увеличения нагрузки на ram")
	t.Description("меняем нагрузку на оперативную память в приложении моке")
	ises.prepareConfig(t)

	ises.changeTestAppCpuLoad(t)
}

func (ises *T09Suite) changeTestAppCpuLoad(t provider.T) {
	t.WithNewStep("Выдаем назгрузку на приложение", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.AppName, ises.testConfig.LoadNamespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			changeLoadStep := allure.NewSimpleStep("Меняем memory нагрузку этому поду")
			values := map[string]int{"memoryPercent": ises.testConfig.Memload * 2350}

			jsonValue, _ := json.Marshal(values)
			string, err := ises.k8sRepo.PodSendPostRequest(podName, ises.testConfig.LoadNamespace, "http://localhost:5000/memory", "5000:3000", kubernetes.Data{JsonData: jsonValue})
			if err != nil {
				fmt.Println(string)
				t.Fatal(err.Error())
			}
			changeLoadStep.WithParent(getPodStep)
			time.Sleep(time.Duration(ises.testConfig.Timeout) * time.Second)
			logStr := fmt.Sprintf(" выдана нагрузка на ram = %v ", ises.testConfig.Memload)
			getPodStep.WithAttachments(allure.NewAttachment("Нагрузка выдана", allure.Text, []byte(logStr)))
		}
	})
}
