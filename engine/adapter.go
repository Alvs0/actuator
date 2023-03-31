package engine

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	_ "github.com/go-sql-driver/mysql"
)

type SqlConfig struct {
	DbName     string `json:"dbName"`
	DbUser     string `json:"dbUser"`
	DbHost     string `json:"dbHost"`
	DbPort     string `json:"dbPort"`
	DbPassword string `json:"dbPassword"`
	Region     string `json:"region"`
}

type SqlAdapter interface {
	Read(query string, input []interface{}, output interface{}) error
	Write(query string, input []interface{}) error
}

type sqlAdapter struct {
	sqlDb *sql.DB
}

func NewSqlAdapter(sqlCfg SqlConfig) SqlAdapter {
	dbEndpoint := fmt.Sprintf("%s:%s", sqlCfg.DbHost, sqlCfg.DbPort)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		sqlCfg.DbUser, sqlCfg.DbPassword, dbEndpoint, sqlCfg.DbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("[SqlAdapter] failed to open database. cause: ", err.Error())
	}

	return &sqlAdapter{
		sqlDb: db,
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
