package impl

import (
	"actuator/engine"
	"fmt"
	"log"
	"strings"
	"time"
)

type SensorQuery interface {
	UpsertSensor(sensorDbs []SensorDb) error
	GetSensors(filter SensorFilter) ([]SensorDb, error)
	DeleteSensor(filter SensorFilter) error
}

type sensorQuery struct {
	sqlAdapter engine.SqlAdapter
}

func NewSensorQuery(sqlAdapter engine.SqlAdapter) SensorQuery {
	return &sensorQuery{
		sqlAdapter: sqlAdapter,
	}
}

type SensorDb struct {
	FirstID     string    `json:"first_id" db:"first_id"`
	SecondID    string    `json:"second_id" db:"second_id"`
	SensorValue string    `json:"sensor_value" db:"sensor_value"`
	SensorType  string    `json:"sensor_type" db:"sensor_type"`
	Timestamp   time.Time `json:"timestamps" db:"timestamps"`
}

type SensorFilter struct {
	FirstID   *string
	SecondID  *string
	Timestamp *time.Time
}

func (s *SensorFilter) ConstructFilter(ignoreEmpty bool) (string, []interface{}, error) {
	var firstIDExist, secondIDExist, timestampExist bool

	var whereClause []string
	var whereInput []interface{}
	if s.FirstID != nil && *s.FirstID != "" {
		firstIDExist = true
		whereClause = append(whereClause, "first_id = ?")
		whereInput = append(whereInput, *s.FirstID)
	}

	if s.SecondID != nil && *s.SecondID != "" {
		secondIDExist = true
		whereClause = append(whereClause, "second_id = ?")
		whereInput = append(whereInput, *s.SecondID)
	}

	if s.Timestamp != nil {
		timestampExist = true
		whereClause = append(whereClause, "timestamps = ?")
		whereInput = append(whereInput, *s.Timestamp)
	}

	if !ignoreEmpty && !firstIDExist && !secondIDExist && !timestampExist {
		return "", nil, fmt.Errorf("empty filter disallowed")
	}

	if ignoreEmpty && !firstIDExist && !secondIDExist && !timestampExist {
		return "", nil, nil
	}

	whereString := strings.Join(whereClause, " AND ")

	return fmt.Sprintf("WHERE %v", whereString), whereInput, nil
}

func (a *sensorQuery) UpsertSensor(sensorDbs []SensorDb) error {
	insertQuery := `INSERT INTO sensor(
							first_id, 
							second_id, 
							sensor_value, 
							sensor_type, 
							timestamps
					) %v AS new(a,b,c,d,e)
					ON DUPLICATE KEY UPDATE 
							sensor_value = c, 
							sensor_type = d,
							timestamps = e`

	var values []string
	var queryInput []interface{}
	for _, sensorDb := range sensorDbs {
		values = append(values, `VALUES(?, ?, ?, ?, ?)`)

		queryInput = append(queryInput, sensorDb.FirstID)
		queryInput = append(queryInput, sensorDb.SecondID)
		queryInput = append(queryInput, sensorDb.SensorValue)
		queryInput = append(queryInput, sensorDb.SensorType)
		queryInput = append(queryInput, sensorDb.Timestamp)
	}

	valueString := strings.Join(values, ", ")

	finalQuery := fmt.Sprintf(insertQuery, valueString)

	if err := a.sqlAdapter.Write(finalQuery, queryInput); err != nil {
		log.Printf("[UpsertSensor] error upserting sensor. cause: %v\n", err.Error())
		return err
	}

	return nil
}

func (a *sensorQuery) GetSensors(filter SensorFilter) ([]SensorDb, error) {
	whereString, queryInput, err := filter.ConstructFilter(true)
	if err != nil {
		return nil, err
	}

	getQuery := fmt.Sprintf(`SELECT * FROM sensor %v`, whereString)

	var sensorDbs []SensorDb
	if err := a.sqlAdapter.Read(getQuery, queryInput, &sensorDbs); err != nil {
		return nil, err
	}

	return sensorDbs, nil
}

func (a *sensorQuery) DeleteSensor(filter SensorFilter) error {
	whereString, queryInput, err := filter.ConstructFilter(false)
	if err != nil {
		return err
	}

	deleteQuery := fmt.Sprintf(`DELETE FROM sensor %v`, whereString)
	if err := a.sqlAdapter.Write(deleteQuery, queryInput); err != nil {
		return err
	}

	return nil
}
