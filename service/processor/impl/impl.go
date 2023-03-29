package impl

import (
	processor "actuator/service/processor/api"
)

type processorService struct {
	processor.UnimplementedProcessorServer
}
