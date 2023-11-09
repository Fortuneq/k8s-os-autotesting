package t05

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"k8s.io/apimachinery/pkg/util/yaml"
	"time"
)

type T05Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T05Config
}

func (ises *T05Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T05Suite) New() model.RunnableTest {
	return &T05Suite{}
}

func (ises *T05Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T05Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T05Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T05Config](b)
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

func (ises *T05Suite) Test(t provider.T) {
	t.Title("Изменение нагрузки в mock приложении для увеличения нагрузки на cpu")
	t.Description("меняем нагрузку на процессор в приложении моке ")
	ises.prepareConfig(t)

	ises.changeTestAppCpuLoad(t)
}

func (ises *T05Suite) changeTestAppCpuLoad(t provider.T) {
	t.WithNewStep("Выдаем назгрузку на приложение", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.AppName, ises.testConfig.LoadNamespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			changeLoadStep := allure.NewSimpleStep("Меняем cpu нагрузку этому поду")
			values := map[string]int{"cpuPercent": ises.testConfig.Cpuload}

			jsonValue, _ := json.Marshal(values)
			string, err := ises.k8sRepo.PodSendPostRequest(podName, ises.testConfig.LoadNamespace, "http://localhost:3000/", "3000:3000", kubernetes.Data{JsonData: jsonValue})
			if err != nil {
				fmt.Println(string)
				t.Fatal(err.Error())
			}
			changeLoadStep.WithParent(getPodStep)
			time.Sleep(time.Duration(ises.testConfig.Timeout) * time.Second)
			logStr := fmt.Sprintf(" выдана нагрузка на cpu = %v ", ises.testConfig.Cpuload)
			getPodStep.WithAttachments(allure.NewAttachment("Нагрузка выдана", allure.Text, []byte(logStr)))
		}
	})
}
