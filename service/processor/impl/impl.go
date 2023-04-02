package impl

import (
	processor "actuator/service/processor/api"
)

type ProcessorService struct {
	ProcessorQuery ProcessorQuery

	processor.UnimplementedProcessorServer
}
