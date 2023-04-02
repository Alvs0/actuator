package impl

import (
	"context"
	processor "github.com/Alvs0/actuator/service/processor/api"
	"github.com/golang/protobuf/ptypes/empty"
)

func (p *ProcessorService) Delete(ctx context.Context, req *processor.SensorFilter) (res *empty.Empty, err error) {
	if err := p.ProcessorQuery.DeleteSensor(constructSensorFilter(req)); err != nil {
		return new(empty.Empty), err
	}

	res = &empty.Empty{}
	return
}
