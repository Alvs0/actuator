package impl

import (
	processor "actuator/service/processor/api"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"strconv"
)

func (p *ProcessorService) Edit(ctx context.Context, req *processor.EditRequest) (res *empty.Empty, err error) {
	var sensorDbs []SensorDbUpsertSpec
	for _, sensorObj := range req.GetSensors() {
		sensorDbs = append(sensorDbs, SensorDbUpsertSpec{
			FirstID:     sensorObj.GetId1(),
			SecondID:    strconv.Itoa(int(sensorObj.GetId2())),
			SensorValue: sensorObj.GetSensorValue(),
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
