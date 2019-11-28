package laas

import (
	"go.uber.org/zap"
	"time"

	"go.uber.org/cadence/workflow"
)

func init() {
	workflow.Register(SimpleWorkflow)
}

func SimpleWorkflow(ctx workflow.Context, value string) error {
	ao := workflow.ActivityOptions{
		TaskList:               "Source",
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	future := workflow.ExecuteActivity(ctx, SimpleActivity, value)
	var result string
	if err := future.Get(ctx, &result); err != nil {
		return err
	}
	workflow.GetLogger(ctx).Info("Done", zap.String("result", result))
	return nil
}