package main

import (
  "context"
  "fmt"
  "github.com/google/uuid"
  "go.uber.org/cadence/client"
  "go.uber.org/zap"
  "time"
)

func TriggerFunction() {
  logger, err := BuildLogger()
  if err != nil {
    panic(err.Error())
  }
  startWorkflow(logger)
}

func startWorkflow(logger *zap.Logger) {
  id := fmt.Sprintf("helloworld_%d", uuid.New())
  workflowOptions := client.StartWorkflowOptions{
    ID:                              id,
    TaskList:                        TaskListName,
    ExecutionStartToCloseTimeout:    time.Minute,
    DecisionTaskStartToCloseTimeout: time.Minute,
  }

  workflowClient, err := BuildCadenceClient(logger)
  workflowExecution, err := workflowClient.StartWorkflow(context.Background(),
    workflowOptions, SimpleWorkflow, "value to Simpleworkflow")
  if err != nil {
    logger.Error(err.Error())
    panic("Failed to start worker")
  }

  logger.Info("Started Workflow.",
    zap.String("Workflow ID: ", workflowExecution.ID),
    zap.String("Workflow Run ID: ", workflowExecution.RunID))
}
