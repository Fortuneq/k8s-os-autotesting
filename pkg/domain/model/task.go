package model

type Task struct {
	Jobs       []JobConfig
	Tests      []TestConfig
	Finalizers []JobConfig
}
