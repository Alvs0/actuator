package impl

import (
	processor "actuator/service/processor/api"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

const (
	TimeLayout = "2006-01-02T15:04:05Z"
)

func (p *ProcessorService) Get(ctx context.Context, req *processor.SensorFilterAndPagination) (res *processor.SensorResponse, err error) {
	sensorDbs, err := p.ProcessorQuery.GetSensors(constructSensorFilterAndPagination(req))
	if err != nil {
		return new(processor.SensorResponse), err
	}

	var sensors []*processor.Sensor
	for _, sensorDb := range sensorDbs {
		secondIDInt, err := strconv.Atoi(sensorDb.SecondID)
		if err != nil {
			continue
		}

		timestamp, err := time.Parse(TimeLayout, sensorDb.Timestamp)
		if err != nil {
			continue
		}

		sensor := processor.Sensor{
			SensorValue: sensorDb.SensorValue,
			SensorType:  sensorDb.SensorType,
			Id1:         sensorDb.FirstID,
			Id2:         int32(secondIDInt),
			Timestamp:   timestamppb.New(timestamp),
		}

		sensors = append(sensors, &sensor)
	}

	res = &processor.SensorResponse{
		Sensors: sensors,
	}

	return
}

func constructSensorFilter(filterReq *processor.SensorFilter) SensorFilter {
	var firstID, secondID *string
	var startTimestamp, endTimestamp *time.Time

	if filterReq.GetId1() != nil || filterReq.GetId1().GetData() != "" {
		firstIDData := filterReq.GetId1().GetData()
		firstID = &firstIDData
	}

	if filterReq.GetId2() != nil || filterReq.GetId2().GetData() != 0 {
		secondIDData := strconv.Itoa(int(filterReq.GetId2().GetData()))
		secondID = &secondIDData
	}

	if filterReq.GetStartTimestamp() != nil {
		startTimestampData := filterReq.GetStartTimestamp().GetData().AsTime()
		startTimestamp = &startTimestampData
	}

	if filterReq.GetEndTimestamp() != nil {
		endTimestampData := filterReq.GetEndTimestamp().GetData().AsTime()
		endTimestamp = &endTimestampData
	}

	return SensorFilter{
		FirstID:        firstID,
		SecondID:       secondID,
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}
}

func constructSensorFilterAndPagination(req *processor.SensorFilterAndPagination) (SensorFilter, SensorPagination) {
	pageNumber := req.GetSensorPagination().GetPageNumbers()
	itemPerPage := req.GetSensorPagination().GetItemPerPage()

	return constructSensorFilter(req.GetSensorFilter()), SensorPagination{
		PageNumber:  &pageNumber,
		ItemPerPage: &itemPerPage,
	}
}
