package impl

import (
	processor "actuator/service/processor/api"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

func (p *ProcessorService) Get(ctx context.Context, req *processor.SensorFilterAndPagination) (res *processor.SensorResponse, err error) {
	firstID := req.SensorFilter.GetId1()
	secondID := strconv.Itoa(int(req.SensorFilter.GetId2()))
	timestamp := req.SensorFilter.GetTimestamp().AsTime()

	sensorDbs, err := p.SensorQuery.GetSensors(SensorFilter{
		FirstID:   &firstID,
		SecondID:  &secondID,
		Timestamp: &timestamp,
	})
	if err != nil {
		return nil, err
	}

	var sensors []processor.Sensor
	for _, sensorDb := range sensorDbs {
		secondIDInt, err := strconv.Atoi(sensorDb.SecondID)
		if err != nil {
			continue
		}

		sensorValueFloat64, err := strconv.ParseFloat(sensorDb.SensorValue, 10)
		if err != nil {
			continue
		}

		sensors = append(sensors, processor.Sensor{
			SensorValue:          float32(sensorValueFloat64),
			SensorType:           sensorDb.SensorType,
			Id1:                  sensorDb.FirstID,
			Id2:                  int32(secondIDInt),
			Timestamp:            timestamppb.New(sensorDb.Timestamp),
		})
	}

	res = &processor.SensorResponse{
		Sensors:              nil,
	}

	return nil, nil
}
