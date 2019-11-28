package cadence

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/uber-go/tally"
	"go.uber.org/cadence/client"
	"go.uber.org/zap"
	//"go.uber.org/zap/zapcore"
	"time"
)

//var HostPort = "10.2.102.76:8088"
//var Domain = "LAAS"
//var TaskListName = "Source"
//var ClientName = "source-worker"
//var CadenceService = "cadence-frontend"


func trigger() {
	startWorkflow(BuildLogger())
}
//
//func buildLogger() *zap.Logger {
//	config := zap.NewDevelopmentConfig()
//	config.Level.SetLevel(zapcore.InfoLevel)
//
//	var err error
//	logger, err := config.Build()
//	if err != nil {
//		panic("Failed to setup logger")
//	}
//
//	return logger
//}

//func buildServiceClient() (workflowserviceclient.Interface, error) {
//	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
//	if err != nil {
//		panic("Failed to setup tchannel")
//	}
//	dispatcher := yarpc.NewDispatcher(yarpc.Config{
//		Name: ClientName,
//		Outbounds: yarpc.Outbounds{
//			CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
//		},
//	})
//	if err := dispatcher.Start(); err != nil {
//		panic("Failed to start dispatcher")
//	}
//
//	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService)), nil
//}

func buildCadenceClient() (client.Client, error) {
	service := BuildServiceClient()
	clientOptions := client.Options{
		Identity: "",
		MetricsScope: tally.NoopScope,
	}
	return client.NewClient(service, Domain, &clientOptions), nil
}

func startWorkflow(logger *zap.Logger) {
	id := fmt.Sprintf("helloworld_%d", uuid.New())
	workflowOptions := client.StartWorkflowOptions{
		ID:                              id,
		TaskList:                        TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	workflowClient, err := buildCadenceClient()
	workflowExecution, err := workflowClient.StartWorkflow(context.Background(), workflowOptions, SimpleWorkflow)
	if err != nil {
		panic("Failed to start worker")
	}

	logger.Info("Started Workflow.",
		zap.String("Workflow ID: ", workflowExecution.ID),
		zap.String("Workflow Run ID: ", workflowExecution.RunID))
}