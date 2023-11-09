package model

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

type RunnableTest interface {
	runner.TestSuite
	Test(t provider.T)
	SetTestConfig(config TestConfig)
	New() RunnableTest
}

type TestConfig struct {
	Name   string
	Type   string
	Params any
}
