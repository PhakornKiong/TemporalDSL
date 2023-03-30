package dsl

import (
	"context"
	"fmt"

	"encore.dev/rlog"
	"go.temporal.io/sdk/activity"
)

type SampleActivities struct {
}

func (a *SampleActivities) SampleActivity1(ctx context.Context, input interface{}) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	rlog.Info(fmt.Sprintf("Run %s with input %v \n", name, input))
	return fmt.Sprint(input), nil
}

func (a *SampleActivities) SampleActivity2(ctx context.Context, input interface{}) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	rlog.Info(fmt.Sprintf("Run %s with input %v \n", name, input))
	return "Result_" + name, nil
}

func (a *SampleActivities) SampleActivity3(ctx context.Context, input interface{}) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	rlog.Info(fmt.Sprintf("Run %s with input %v \n", name, input))
	return "Result_" + name, nil
}

func (a *SampleActivities) SampleActivity4(ctx context.Context, input interface{}) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	rlog.Info(fmt.Sprintf("Run %s with input %v \n", name, input))
	return "Result_" + name, nil
}
