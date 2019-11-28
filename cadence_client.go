package main

import (
  "fmt"
  "github.com/uber-go/tally"
  "go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
  "go.uber.org/cadence/client"
  "go.uber.org/yarpc"
  "go.uber.org/yarpc/transport/tchannel"
  "go.uber.org/zap"
  "go.uber.org/zap/zapcore"
)

var HostPort = "10.2.145.115:7933"
var Domain = "laas_cadence"
var TaskListName = "Source"
var ClientName = "simpleworker"
var CadenceService = "cadence-frontend"

//var HostPort = "10.2.102.76:7933"
//var Domain = "LAAS"
//var TaskListName = "Source"
//var ClientName = "simpleworker"
//var CadenceService = "cadence-frontend"

func BuildLogger() (*zap.Logger, error) {
  config := zap.NewDevelopmentConfig()
  config.Level.SetLevel(zapcore.InfoLevel)

  var err error
  logger, err := config.Build()
  if err != nil {
    fmt.Errorf(err.Error())
    return nil, err
  }

  return logger, nil
}

func BuildServiceClient(logger *zap.Logger) (workflowserviceclient.Interface, error) {
  ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
  if err != nil {
    logger.Error(err.Error())
    return nil, err
  }
  dispatcher := yarpc.NewDispatcher(yarpc.Config{
    Name: ClientName,
    Outbounds: yarpc.Outbounds{
      CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
    },
  })
  if err := dispatcher.Start(); err != nil {
    logger.Error(err.Error())
    return nil, err
  }

  return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService)), nil
}

func BuildCadenceClient(logger *zap.Logger) (client.Client, error) {
  service, err := BuildServiceClient(logger)
  if err != nil {
    logger.Error(err.Error())
    return nil, err
  }
  clientOptions := client.Options{
    Identity:     "",
    MetricsScope: tally.NoopScope,
  }
  return client.NewClient(service, Domain, &clientOptions), nil
}
