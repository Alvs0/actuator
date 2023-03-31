package impl

import (
	processor "actuator/service/processor/api"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"strconv"
)

func (p *ProcessorService) Edit(ctx context.Context, req *processor.EditRequest) (res *empty.Empty, err error) {
	var sensorDbs []SensorDb
	for _, sensorObj := range req.GetSensors() {
		sensorDbs = append(sensorDbs, SensorDb{
			FirstID:     sensorObj.GetId1(),
			SecondID:    strconv.Itoa(int(sensorObj.GetId2())),
			SensorValue: fmt.Sprintf("%f", sensorObj.GetSensorValue()),
			SensorType:  sensorObj.GetSensorType(),
			Timestamp:   sensorObj.GetTimestamp().AsTime(),
		})
	}

	if err := p.SensorQuery.UpsertSensor(sensorDbs); err != nil {
		return nil, err
	}

	res = &empty.Empty{}
	return nil, nil
}
