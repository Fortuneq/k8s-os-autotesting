package utils

import (
	"log"
	"time"
)

// Функция ретрая
type RetryFn func() (string, error)

// Функция фильтра
type FilterFn func(string) bool

// Объект для запуска ретрая
type RetryWithResponse struct {
	// Количество попыток
	N int
	// Время перерыва между попытками
	Sleep time.Duration
	// Функция ретрая
	Fn RetryFn
	// Функция фильтра
	Filter FilterFn
}

// todo перейти на ретрай кубера
// Запустить ретрай
func (r RetryWithResponse) Start() (response string, err error) {
	for i := 0; i < r.N; i++ {
		if response, err = r.Fn(); err == nil {
			// Если задан фильтр на ответ, то проверяем по нему
			if r.Filter != nil {
				// Прошел по фильтру - выходим
				if r.Filter(response) {
					break
				}
			} else {
				// Выходим, если нет фильтра и запрос прошел успешно
				break
			}
		}
		log.Println("Sleep", r.Sleep)
		time.Sleep(r.Sleep)
	}

	return
}
