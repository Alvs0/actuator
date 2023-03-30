package impl

import (
	generator "actuator/service/generator/api"
	processor "actuator/service/processor/api"
	"os"
	"sync"
)

type GeneratorService struct {
	Worker Worker
	ProcessorService processor.ProcessorClient

	generator.UnimplementedGeneratorServer
}

type Worker struct {
	Workers     int
	MessageChan chan struct{}
	SignalChan  chan os.Signal
	waitGroup   sync.WaitGroup
}

func NewWorker(numOfWorkers int, batch int) Worker {
	return Worker{
		Workers:     numOfWorkers,
		MessageChan: make(chan struct{}, batch),
		SignalChan:  make(chan os.Signal, 1),
		waitGroup:   sync.WaitGroup{},
	}
}