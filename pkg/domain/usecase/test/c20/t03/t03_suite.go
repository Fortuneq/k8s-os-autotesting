package t03

import (
	"context"
	"fmt"
	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"k8s.io/apimachinery/pkg/util/yaml"
	"log"
	"math"
)

type T03Suite struct {
	suite.Suite
	k8sRepo          kubernetes.K8SRepo
	istioRepo        istio.IstioRepo
	commonTestConfig model.TestConfig
	testConfig       T03Config
}

func (ises *T03Suite) SetTestConfig(config model.TestConfig) {
	ises.commonTestConfig = config
}

func (ises T03Suite) New() model.RunnableTest {
	return &T03Suite{}
}

func (ises *T03Suite) SetK8sRepo(repo kubernetes.K8SRepo) {
	ises.k8sRepo = repo
}

func (ises *T03Suite) SetIstioRepo(repo istio.IstioRepo) {
	ises.istioRepo = repo
}

func (ises *T03Suite) prepareConfig(t provider.T) {
	if b, err := yaml.Marshal(ises.commonTestConfig.Params); err == nil {
		params, err := config.ReadConfigFromByteArray[T03Config](b)
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

func (ises *T03Suite) Test(t provider.T) {
	t.Title("Тестирование количетсва подов при скейлинге по утилизации cpu")
	t.Description("Проверяем количетсво подов приложения по текущим cpu параметрам")
	ises.prepareConfig(t)

	ises.checkTestAppReady(t)
}

func (ises *T03Suite) checkTestAppReady(t provider.T) {
	t.WithNewStep("Проверяем количетсво подов приложения по текущим cpu параметрам", func(ctx provider.StepCtx) {
		if ises.testConfig.Workload == params.DeploymentWorkload {
			getPodStep := allure.NewSimpleStep("Ищем podы тестового приложения по названию deployment")
			getPodStep.WithParent(ctx.CurrentStep())
			podName, err := ises.k8sRepo.GetAllDeploymentPodNames(context.TODO(), ises.testConfig.AppName, ises.testConfig.TestNamespace)
			if err != nil {
				t.Fatal(err.Error())
			}
			cpuUsage := ises.k8sRepo.GetCurrentUsageMetrics(ises.testConfig.AppName, ises.testConfig.TestNamespace)
			if cpuUsage.Cpu == 0 {
				t.Fatal(err.Error())
			}
			limits := ises.k8sRepo.GetLimits(ises.testConfig.AppName, ises.testConfig.TestNamespace)
			cpuCoef := 0.0
			if limits.Cpu != 0.0 {
				cpuCoef = cpuUsage.Cpu / limits.Cpu
			}
			newReplicas := 1.0
			if cpuCoef != 0.0 {
				newReplicas = math.Max(newReplicas, math.Ceil(float64(cpuUsage.PodsCount)*cpuCoef/(float64(ises.testConfig.Cpuload)/100)))
				log.Printf("посчитанное количество реплик %f cpuPercentage: %f", newReplicas, float64(ises.testConfig.Cpuload))
			}
			if int(newReplicas) != len(podName) {
				t.Fatal(fmt.Errorf("посчитанное количество реплик %v сколько в кластере  %v", int(newReplicas), len(podName)))
			}
			logStr := fmt.Sprintf(" %v реплик поднялись, status code = 200", len(podName))
			getPodStep.WithAttachments(allure.NewAttachment("Количество реплик совпадает с необходимым", allure.Text, []byte(logStr)))
		}
	})
}
