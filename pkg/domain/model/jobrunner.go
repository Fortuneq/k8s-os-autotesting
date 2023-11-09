package model

import "context"

type JobRunner func(context.Context, JobConfig) error
