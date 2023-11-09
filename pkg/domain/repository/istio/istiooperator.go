package istio

import (
	"bytes"
	"context"
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"text/template"
	"time"

	operatorv1alpha1 "istio.io/api/operator/v1alpha1"
	iopv1alpha1 "istio.io/istio/operator/pkg/apis/istio/v1alpha1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"sbet-tech.com/synapse/istio-se/allure/pkg/domain/usecase/job/params"
	"sbet-tech.com/synapse/istio-se/allure/pkg/utils"
)

const (
	defaultWatchTimeout = time.Minute * 5
)

func (ir *istioRepo) DeleteIstioOperator(params params.IstioOperatorDeleteParams) error {
	return retry.RetryOnConflict(wait.Backoff{
		Steps:    10,
		Duration: 5 * time.Second,
		Factor:   1.5,
		Jitter:   0.1,
	}, func() error {
		var iop iopv1alpha1.IstioOperator

		if err := ir.k8sRuntimeClient.Get(context.TODO(), client.ObjectKey{
			Name:      params.Name,
			Namespace: params.Namespace},
			&iop, &client.GetOptions{},
		); err != nil {
			if kerrors.IsNotFound(err) {
				return nil
			}
			return err
		}

		return ir.k8sRuntimeClient.Delete(context.TODO(), &iop, &client.DeleteOptions{})
	})
}

// Создаем IstioOperator и ожидаем в течение некоторого таймаута пока его Status не станет HEALTHY
func (ir *istioRepo) CreateIstioOperator(params params.IstioOperatorInstallParams) error {
	var (
		iop iopv1alpha1.IstioOperator
		obj []byte
	)

	if params.Resource != "" {
		// Нам задали ресурс в параметрах
		obj = []byte(params.Resource)
	} else {
		var err error
		if obj, err = ioutil.ReadFile(params.Path); err != nil {
			return err
		}
	}

	if params.Values != "" {
		var (
			temp *template.Template
			err  error
		)
		if temp, err = template.New("new").Parse(string(obj)); err != nil {
			return err
		}

		var templateValues any
		err = yaml.Unmarshal([]byte(params.Values), &templateValues)
		if err != nil {
			return err
		}

		buf := &bytes.Buffer{}
		err = temp.Execute(buf, templateValues)
		if err != nil {
			return err
		}

		obj = buf.Bytes()
	}

	// YAML -> K8S ресурс
	if err := utils.ReadYamlToObject(obj, &iop); err != nil {
		return err
	}

	// Задаем таймаут
	var t time.Duration
	var err error
	if params.Timeout != "" {
		if t, err = time.ParseDuration(params.Timeout); err != nil {
			return err
		}
	} else {
		t = defaultWatchTimeout
	}
	// Создаем IstioOperator
	if err := ir.k8sRuntimeClient.Get(context.TODO(), types.NamespacedName{Namespace: iop.Namespace, Name: iop.Name}, &iop); err == nil {
		// Уже есть на кластере
		log.Printf("Update resource %s in namespace %s", iop.Name, iop.Namespace)
		if err := ir.k8sRuntimeClient.Update(context.TODO(), &iop, &client.UpdateOptions{}); err != nil {
			return err
		}
		log.Println("Start watching", iop.Name, "in", iop.Namespace, "for", t)
		ctx, cancel := context.WithTimeout(context.TODO(), t)
		defer cancel()
		return ir.watchForStatusChange(ctx, iop.Namespace, iop.Name)
	} else {
		// Нет на кластере, создаем
		log.Printf("Create resource %s in namespace %s", iop.Name, iop.Namespace)
		if err := ir.k8sRuntimeClient.Create(context.TODO(), &iop, &client.CreateOptions{}); err != nil {
			return err
		} else {
			log.Println("Start watching", iop.Name, "in", iop.Namespace, "for", t)
			ctx, cancel := context.WithTimeout(context.TODO(), t)
			defer cancel()
			return ir.watchForStatusChange(ctx, iop.Namespace, iop.Name)
		}
	}
}

// Чекаем Status IstioOperator пока не станет HEALTHY или таймаут не пикнет
func (ir *istioRepo) watchForStatusChange(ctx context.Context, ns, name string) error {
	// Запускаем Watch
	if w, err := ir.k8sRuntimeClient.Watch(ctx, &iopv1alpha1.IstioOperatorList{}, &client.ListOptions{
		Namespace:     ns,                                                 // Смотрим в конкретном проекте
		FieldSelector: fields.OneTermEqualSelector("metadata.name", name), // Смотрим за конкретным ресурсом
	}); err != nil {
		return err
	} else {
		defer w.Stop()

		for {
			select {
			// Ожидаем, мб Wath выдаст нам изменения
			case e, ok := <-w.ResultChan():
				if !ok {
					return errors.New("chan is closed")
				}
				// Чекаем, что словили IstioOperator
				if v, ok := e.Object.(*iopv1alpha1.IstioOperator); !ok {
					continue
				} else {
					if v.Status != nil {
						// Статус все еще не в нужном значении
						if v.Status.Status != operatorv1alpha1.InstallStatus_HEALTHY {
							log.Println(v.Status.Status)
						} else {
							// Повезло повезло
							log.Println("SUCCESS!", v.Status.Status)
							return nil
						}
					}
				}
			// На всякий случай, раз в 10 секунд делаем Get и чекаем его статус, на тот кейс если будут проблемы с Watch
			case <-time.Tick(time.Second * 10):
				log.Println("Watch Tick")
				if err := getIstioOperatorAndCheckStatus(ns, name); err != nil {
					log.Println("Status field in IstioOperator is still not HEALTHY")
				} else {
					return nil
				}
			// Время вышло
			case <-ctx.Done():
				// Перед выходом на посошок чекнем. что статус поменялся. Мож повезет
				if err := getIstioOperatorAndCheckStatus(ns, name); err != nil {
					return errors.New("deadline")
				}
				return nil
			}
		}
	}
}

// Получаем ресурс и чекаем его статус
func getIstioOperatorAndCheckStatus(ns, name string) error {
	var iop iopv1alpha1.IstioOperator
	if err := ir.k8sRuntimeClient.Get(context.TODO(), types.NamespacedName{Namespace: ns, Name: name}, &iop, &client.GetOptions{}); err != nil {
		return err
	} else if iop.Status != nil && iop.Status.Status == operatorv1alpha1.InstallStatus_HEALTHY {
		log.Println("SUCCESS!", iop.Status.Status)
		return nil
	}

	return errors.New("IstioOperator Status field is not HEALTHY")
}

func (istioRepo *istioRepo) GetIstioOperator(name, namespace string) (*iopv1alpha1.IstioOperator, error) {
	var iop iopv1alpha1.IstioOperator

	if err := istioRepo.k8sRuntimeClient.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, &iop, &client.GetOptions{}); err != nil {
		return &iopv1alpha1.IstioOperator{}, err
	} else if iop.Status != nil && iop.Status.Status == operatorv1alpha1.InstallStatus_HEALTHY {
		//TODO мы уверены, что нам это нужно?
		log.Println(iop.Spec.Values.Fields["pilot"].GetStructValue().Fields["namespaceWideValidation"].GetBoolValue())
		return &iop, nil
	}
	return &iopv1alpha1.IstioOperator{}, nil
}

func (istioRepo *istioRepo) UpdateIstioOperator(iop *iopv1alpha1.IstioOperator) error {
	return istioRepo.k8sRuntimeClient.Update(context.TODO(), iop)
}

func (istioRepo *istioRepo) PatchIstioOperator(iop *iopv1alpha1.IstioOperator, patch []byte) error {
	return istioRepo.k8sRuntimeClient.Patch(context.TODO(), iop, client.RawPatch(types.MergePatchType, patch))
}
