package impl

import (
	processor "actuator/service/processor/api"
	"fmt"
	"io"
)

func (p *processorService) Process(stream processor.Processor_ProcessServer) error {
	// ToDo: Add Authorization

	var total int32
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return stream.Send(&processor.ProcessResponse{
				Total: total,
			})
		}

		fmt.Println(value)

		total += 1
	}
	return nil
}
