package impl

import (
	"actuator/engine"
	"fmt"
	"log"
	"strings"
	"time"
)

type SensorQuery interface {
	UpsertSensor(sensorDb SensorDb) error
	GetSensors(ids []string) ([]SensorDb, error)
	DeleteSensor(id string) error
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

func (a *sensorQuery) UpsertSensor(sensorDb SensorDb) error {
	insertQuery := `INSERT INTO sensor(
							first_id, 
							second_id, 
							sensor_value, 
							sensor_type, 
							timestamps
					) VALUES(?, ?, ?, ?, ?) AS new(a,b,c,d,e)
					ON DUPLICATE KEY UPDATE 
							sensor_value = c, 
							sensor_type = d,
							timestamps = e`

	queryInput := []interface{}{
		sensorDb.FirstID,
		sensorDb.SecondID,
		sensorDb.SensorValue,
		sensorDb.SensorType,
		sensorDb.Timestamp,
	}

	if err := a.sqlAdapter.Write(insertQuery, queryInput); err != nil {
		log.Printf("[UpsertSensor] error upserting sensor. cause: %v\n", err.Error())
		return err
	}

	return nil
}

func (a *sensorQuery) GetSensors(ids []string) ([]SensorDb, error) {
	getQuery := fmt.Sprintf(`SELECT * FROM "public".sensor WHERE id IN ('%s')`, strings.Join(ids, "', '"))

	var sensorDbs []SensorDb
	if err := a.sqlAdapter.Read(getQuery, nil, &sensorDbs); err != nil {
		return nil, err
	}

	return sensorDbs, nil
}

func (a *sensorQuery) DeleteSensor(id string) error {
	insertQuery := `DELETE FROM "public".sensor WHERE first_id = $1`

	queryInput := []interface{}{
		id,
	}

	if err := a.sqlAdapter.Write(insertQuery, queryInput); err != nil {
		return err
	}

	return nil
}
