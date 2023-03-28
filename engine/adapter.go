package engine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"log"
	"reflect"

	_ "github.com/lib/pq"
)

const (
	MySQLDriverName       = "mysql"
	MySQLConnStringFormat = "%v:%v@%v/%v"
)

type SqlConfig struct {
	HostUrl      string `json:"HostUrl"`
	DatabaseName string `json:"DatabaseName"`
	Username     string `json:"Username"`
	Password     string `json:"Password"`
	FullUrl      string `json:"FullUrl"`
}

func (s SqlConfig) GetSqlConnectionSpecification() (driverName string, datasourceName string) {
	driverName = MySQLDriverName
	datasourceName = fmt.Sprintf(MySQLConnStringFormat,
		s.Username,
		s.Password,
		s.HostUrl,
		s.DatabaseName,
	)

	pqUrl, err := pq.ParseURL(s.FullUrl)
	if err != nil {
		log.Fatalln(fmt.Errorf("error parsing database url. cause: %v", err.Error()))
	}

	datasourceName = pqUrl

	return
}

type SqlAdapter interface {
	Read(query string, input []interface{}, output interface{}) error
	Write(query string, input []interface{}) error
}

type sqlAdapter struct {
	sqlDb *sql.DB
}

func NewSqlAdapter(config SqlConfig) SqlAdapter {
	sqlDb, err := sql.Open(config.GetSqlConnectionSpecification())
	if err != nil {
		log.Fatalln(fmt.Errorf("error instantiating database. cause: %v", err.Error()))
	}

	sqlDb.SetMaxIdleConns(10)
	sqlDb.SetMaxOpenConns(10)

	return &sqlAdapter{
		sqlDb: sqlDb,
	}
}

func (s *sqlAdapter) Read(query string, input []interface{}, output interface{}) error {
	res, err := s.sqlDb.Query(query, input...)
	defer res.Close()

	if err != nil {
		return err
	}

	outputType := reflect.TypeOf(output).Elem()
	if outputType.Kind() == reflect.Slice && outputType.Elem().Kind() != reflect.Struct {
		outputPointer := reflect.ValueOf(output)
		outputValue := outputPointer.Elem()
		for res.Next() {
			singularType := reflect.New(outputType.Elem())
			singularInterface := singularType.Interface()
			if err := res.Scan(singularInterface); err != nil {
				return err
			}

			outputValue.Set(reflect.Append(outputValue, reflect.Indirect(reflect.ValueOf(singularInterface))))
		}

		return nil
	} else if outputType.Kind() == reflect.Slice && outputType.Elem().Kind() == reflect.Struct {
		columns, err := res.Columns()
		if err != nil {
			return err
		}

		var allMaps []map[string]interface{}
		for res.Next() {
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(columns))
			for i, _ := range values {
				pointers[i] = &values[i]
			}

			if err := res.Scan(pointers...); err != nil {
				return err
			}

			resultMap := make(map[string]interface{})
			for i, val := range values {
				resultMap[columns[i]] = val
			}

			allMaps = append(allMaps, resultMap)
		}

		allMapsBytes, err := json.Marshal(allMaps)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(allMapsBytes, output); err != nil {
			return err
		}

		return nil
	} else if outputType.Kind() == reflect.Struct {
		columns, err := res.Columns()
		if err != nil {
			return err
		}

		var allMaps map[string]interface{}
		for res.Next() {
			values := make([]interface{}, len(columns))
			pointers := make([]interface{}, len(columns))
			for i, _ := range values {
				pointers[i] = &values[i]
			}

			if err := res.Scan(pointers...); err != nil {
				return err
			}

			resultMap := make(map[string]interface{})
			for i, val := range values {
				resultMap[columns[i]] = val
			}

			allMaps = resultMap
		}

		allMapsBytes, err := json.Marshal(allMaps)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(allMapsBytes, output); err != nil {
			return err
		}

		return nil
	} else {
		if res.Next() {
			if err := res.Scan(output); err != nil {
				return err
			}

			return nil
		}
	}

	fmt.Println(fmt.Sprintf("[SqlAdapter] no result found. query: %v queryInput: %v", query, input))

	return nil
}

func (s *sqlAdapter) Write(query string, input []interface{}) error {
	_, err := s.sqlDb.Exec(query, input...)
	if err != nil {
		return err
	}

	return nil
}
