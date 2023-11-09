package t04

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

type T04Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T04Config
}

func (ises *T04Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T04Suite) New() model.RunnableTest {
	return &T04Suite{}
}

func (ises *T04Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T04Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T04Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T04Config](b)
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

func (ises *T04Suite) Test(t provider.T) {
	t.Title("Тестирование на скейлинг при отключении нагрузки на cpu")
	t.Description("Проверяем, что количество подов даунскейнулось до определенного количетсва")
	ises.prepareConfig(t)

	ises.checkTestAppPodCount(t)
}

func (ises *T04Suite) checkTestAppPodCount(t provider.T) {
	t.WithNewStep("Проверяем Что после даунскейлинга осталось определенное количетсво подов приложения", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetPodName(context.TODO(), ises.testConfig.LoadAppName, ises.testConfig.LoadNamespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			values := map[string]bool{"off": true}

			jsonValue, _ := json.Marshal(values)
			string, err := ises.k8sRepo.PodSendPostRequest(podName, ises.testConfig.LoadNamespace, "http://localhost:3000/", "3000:3000", kubernetes.Data{JsonData: jsonValue})
			if err != nil {
				fmt.Println(string)
				t.Fatal(err.Error())
			}
			time.Sleep(120 * time.Second)
			podNames, err := ises.k8sRepo.GetAllDeploymentPodNames(context.TODO(), ises.testConfig.TestAppName, ises.testConfig.TestNamespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			for i := 0; i < 5; i++ {
				if len(podNames) > ises.testConfig.PodCount {
					time.Sleep(20 * time.Second)
				}
				break
			}
			if len(podNames) > ises.testConfig.PodCount {
				t.Fatal(fmt.Errorf(" количество реплик %v > желаемое %v", len(podNames), ises.testConfig.PodCount))
			}
			logStr := fmt.Sprintf(" %v реплик осталось, status code = 200", len(podNames))
			getPodStep.WithAttachments(allure.NewAttachment("Все поды даунскейнулись", allure.Text, []byte(logStr)))
		}
	})
}
