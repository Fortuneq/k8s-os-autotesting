package test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type K8sRepoSetter interface {
	SetK8sRepo(repo kubernetes.K8SRepo)
}

type IstioRepoSetter interface {
	SetIstioRepo(repo istio.IstioRepo)
}

type K8sTestRunner interface {
	RunTest(testConfig model.TestConfig) error
}

func NewK8sTestRunner(t *testing.T) K8sTestRunner {
	config := utils.KubeRestConfig()
	kr := &k8sTestRunner{
		t:         t,
		k8sRepo:   kubernetes.CreateNewK8SRepo(config),
		istioRepo: istio.CreateNewIstioRepo(config),
	}
	kr.testMap = map[string]model.RunnableTest{
		reflect.TypeOf(synai_t01.T01Suite{}).String(): &synai_t01.T01Suite{},
		reflect.TypeOf(synai_t02.T02Suite{}).String(): &synai_t02.T02Suite{},
		reflect.TypeOf(synai_t03.T03Suite{}).String(): &synai_t03.T03Suite{},
		reflect.TypeOf(synai_t04.T04Suite{}).String(): &synai_t04.T04Suite{},
		reflect.TypeOf(synai_t05.T05Suite{}).String(): &synai_t05.T05Suite{},
		reflect.TypeOf(synai_t06.T06Suite{}).String(): &synai_t06.T06Suite{},
		reflect.TypeOf(synai_t07.T07Suite{}).String(): &synai_t07.T07Suite{},
		reflect.TypeOf(synai_t08.T08Suite{}).String(): &synai_t08.T08Suite{},
		reflect.TypeOf(synai_t09.T09Suite{}).String(): &synai_t09.T09Suite{},
	}
	return kr
}

type k8sTestRunner struct {
	t         *testing.T
	testMap   map[string]model.RunnableTest
	k8sRepo   kubernetes.K8SRepo
	istioRepo istio.IstioRepo
}

func (runner *k8sTestRunner) RunTest(testConfig model.TestConfig) error {
	if test := runner.testMap[testConfig.Type]; test != nil {
		//log.Printf("Config: %+v\n", testConfig)
		test.SetTestConfig(testConfig)
		test.(K8sRepoSetter).SetK8sRepo(runner.k8sRepo)
		test.(IstioRepoSetter).SetIstioRepo(runner.istioRepo)
		suite.RunNamedSuite(runner.t, testConfig.Name, test)
		//suite.RunSuite(runner.t, test)
	} else {
		sb := strings.Builder{}
		for k := range runner.testMap {
			sb.WriteString(k)
			sb.WriteRune('\n')
		}
		//return errors.New(fmt.Sprint("No such test. Available test:", sb.String()))
	}
	return nil
}
