package impl

import (
	processor "actuator/service/processor/api"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (p *ProcessorService) Delete(ctx context.Context, req *processor.SensorFilter) (res *empty.Empty, err error) {
	if err := p.SensorQuery.DeleteSensor(constructSensorFilter(req)); err != nil {
		return new(empty.Empty), err
	}

	res = &empty.Empty{}
	return
}
