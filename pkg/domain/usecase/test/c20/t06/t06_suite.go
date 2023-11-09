package t06

import (
	"context"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/util/yaml"
	"time"
)

type T06Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T06Config
}

func (ises *T06Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T06Suite) New() model.RunnableTest {
	return &T06Suite{}
}

func (ises *T06Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T06Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T06Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T06Config](b)
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

func (ises *T06Suite) Test(t provider.T) {
	t.Title("Тестирование мока монитора")
	t.Description("Проверяем, что количество подов даунскейнулось до определенного количетсва")
	ises.prepareConfig(t)

	ises.requestTestAppMonitor(t)
}

func (ises *T06Suite) requestTestAppMonitor(t provider.T) {
	t.WithNewStep("Отправляем запрос к тестовому приложению моку монитора", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.AppName, ises.testConfig.Namespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			changeLoadStep := allure.NewSimpleStep("Запрос по примеру монитора к  этому поду")
			string, err := ises.k8sRepo.PodSendPostRequest(podName, ises.testConfig.Namespace,
				"http://localhost:3000/coordinator/api/gateway/v1/index/analytical/task/project/ausc/query",
				"3000:3000",
				kubernetes.Data{ApiKey: ises.testConfig.ApiKey, JsonData: []byte(fmt.Sprintf(`{"sqlQuery":%s}`, ises.testConfig.SqlQuery))})
			if err != nil {
				fmt.Println(string)
				t.Fatal(err.Error())
			}
			changeLoadStep.WithParent(getPodStep)
			time.Sleep(60 * time.Second)
			logStr := fmt.Sprintf("Запрос прошел")
			getPodStep.WithAttachments(allure.NewAttachment("Все запросы к поду прошли", allure.Text, []byte(logStr)))
		}
	})
}
