package impl

import (
	processor "actuator/service/processor/api"
)

type ProcessorService struct {
	SensorQuery SensorQuery

	processor.UnimplementedProcessorServer
}
