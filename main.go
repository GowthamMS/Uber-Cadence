package main

import (
  "flag"
)

var workerMode bool

func init() {
  flag.BoolVar(&workerMode, "worker", false, "Start the worker")
}

func main() {
  flag.Parse()
  if workerMode {
    WorkerFunction()
  } else {
    TriggerFunction()
  }
}
