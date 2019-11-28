package main

import (
  "github.com/uber-go/tally"
  "go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
  "go.uber.org/cadence/activity"
  "go.uber.org/cadence/worker"
  "go.uber.org/cadence/workflow"
  "go.uber.org/zap"
)

func init() {
  workflow.Register(SimpleWorkflow)
  activity.Register(SimpleActivity)
}

func WorkerFunction() {
  logger, err := BuildLogger()
  if err != nil {
    panic(err.Error())
  }
  serviceClient, err := BuildServiceClient(logger)
  if err != nil {
    panic(err.Error())
  }
  startWorker(logger, serviceClient)
}

func startWorker(logger *zap.Logger, service workflowserviceclient.Interface) {
  // TaskListName identifies set of client workflows, activities, and workers.
  // It could be your group or client or application name.
  workerOptions := worker.Options{
    Logger:       logger,
    MetricsScope: tally.NewTestScope(TaskListName, map[string]string{}),
  }
  worker := worker.New(
    service,
    Domain,
    TaskListName,
    workerOptions)
  err := worker.Start()
  if err != nil {
    logger.Error(err.Error())
    panic("Failed to start worker")
  }
  logger.Info("Started Worker.", zap.String("worker", TaskListName))
  // For the worker to wait for ever for the workflow to be assigned
  select {}
}
