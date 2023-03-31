package impl

import (
	processor "actuator/service/processor/api"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
)

func (p *ProcessorService) Delete(ctx context.Context, req *processor.SensorFilter) (res *empty.Empty, err error) {
	

	res = &empty.Empty{}
	return nil, nil
}
