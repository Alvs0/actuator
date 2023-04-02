package main

import (
	"actuator/engine"
	generator "actuator/service/generator/api"
	generatorImpl "actuator/service/generator/impl"
	processor "actuator/service/processor/api"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

const (
	ConfigPath = "./config"
)

type Config struct {
	Address          string
	AccountAddress   string
	ProcessorAddress string
	WorkersNumber    int
	Batch            int
}

func main() {
	var cfg Config
	err := engine.LoadConfig("dev", ConfigPath, &cfg)
	if err != nil {
		log.Fatal("[Service] failed to load config cause: ", err.Error())
	}

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalf("failed to listen. cause: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(generatorImpl.Authorize),
	)

	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	processorConn, err := grpc.Dial(cfg.ProcessorAddress, options...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer processorConn.Close()

	processorService := processor.NewProcessorClient(processorConn)
	generatorService := &generatorImpl.GeneratorService{
		Worker:           generatorImpl.NewWorker(cfg.WorkersNumber, cfg.Batch),
		ProcessorService: processorService,
	}

	generator.RegisterGeneratorServer(grpcServer, generatorService)

	fmt.Printf("[Generator] Service ready at %v\n", cfg.Address)
	grpcServer.Serve(listener)
}
