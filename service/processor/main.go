package main

import (
	"fmt"
	"github.com/Alvs0/actuator/engine"
	processor "github.com/Alvs0/actuator/service/processor/api"
	processorImpl "github.com/Alvs0/actuator/service/processor/impl"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	ConfigPath = "./config"
)

type Config struct {
	Address     string
	MySQLConfig engine.SqlConfig
}

func main() {
	var cfg Config
	err := engine.LoadConfig("dev", ConfigPath, &cfg)
	if err != nil {
		log.Fatal("[Service] failed to load config cause: ", err.Error())
	}

	sqlAdapter := engine.NewSqlAdapter(cfg.MySQLConfig)
	sensorQuery := processorImpl.NewProcessorQuery(sqlAdapter)

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalf("failed to listen. cause: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(processorImpl.Authorize),
	)
	processorService := &processorImpl.ProcessorService{
		ProcessorQuery: sensorQuery,
	}

	processor.RegisterProcessorServer(grpcServer, processorService)

	fmt.Printf("[Generator] Service ready at %v\n", cfg.Address)
	grpcServer.Serve(listener)
}
