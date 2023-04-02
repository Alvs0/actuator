package impl

import (
	processor "actuator/service/processor/api"
	"fmt"
	"io"
	"strconv"
)

func (p *ProcessorService) Process(stream processor.Processor_ProcessServer) error {
	// ToDo: Add Authorization

	var total int32
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return stream.Send(&processor.ProcessResponse{
				Total: total,
			})
		}

		if err := p.ProcessorQuery.UpsertSensor([]SensorDbUpsertSpec{
			{
				FirstID:     value.GetId1(),
				SecondID:    strconv.Itoa(int(value.GetId2())),
				SensorValue: value.GetSensorValue(),
				SensorType:  value.GetSensorType(),
				Timestamp:   value.GetTimestamp().AsTime(),
			},
		}); err != nil {
			return err
		}

		fmt.Println("[Processor] Success inserting data")
		total += 1
	}
}
