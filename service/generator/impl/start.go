package impl

import (
	generator "actuator/service/generator/api"
	processor "actuator/service/processor/api"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/jasonlvhit/gocron"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math/rand"
	"time"
)

var randomSensorType = []string{"Thermistor", "IR", "Ultrasonic", "Gyroscope", "Flex", "Humidity", "Smoke"}
var randomCapitalLetter = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (g *GeneratorService) Start(ctx context.Context, req *generator.StartSpec) (res *empty.Empty, err error) {
	stream, err := g.ProcessorService.Process(context.Background())
	if err != nil {
		return nil, err
	}

	go startCron(req, stream)

	res = &empty.Empty{}

	return
}

func startCron(req *generator.StartSpec, stream processor.Processor_ProcessClient) {
	gocron.Every(1).Second().Do(generateRandomSensor, req.NumOfMessagePerSecond, stream)
	<-gocron.Start()
}

func generateRandomSensor(numOfMessagePerSecond int32, stream processor.Processor_ProcessClient) {
	for idx := int32(0); idx < numOfMessagePerSecond; idx++ {
		err := stream.Send(&processor.Sensor{
			SensorValue: rand.Float32(),
			SensorType:  fmt.Sprintf("%v", randomSensorType[rand.Intn(len(randomSensorType))]),
			Id1:         fmt.Sprintf("%v", string(randomCapitalLetter[rand.Intn(len(randomCapitalLetter))])),
			Id2:         int32(rand.Intn(len(randomCapitalLetter))),
			Timestamp:   timestamppb.New(time.Now()),
		})
		if err != nil {
			log.Fatalln("[Error sending value]", err)
		}
	}
}
