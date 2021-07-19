package infoblog

import (
	"fmt"
	"reflect"
	"strings"
)

var entityFields = map[string][]string{}
var entityUpdateFields = map[string][]string{}

func init() {

	entities := map[string]interface{}{
		"users":      User{},
		"likes":      Like{},
		"hashtags":   HashTag{},
		"subscribes": Subscriber{},
		"files":      File{},
		"posts":      PostEntity{},
	}

	for name, entity := range entities {
		t := reflect.TypeOf(entity)

		var allFields []string
		var updateFields []string

		for i := 0; i < t.NumField(); i++ {
			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			field := t.Field(i)

			// Get the field tag value
			tag := field.Tag.Get("db")
			if tag == "" || tag == "-" {
				continue
			}

			tagsString := field.Tag.Get("ops")
			tags := strings.Split(tagsString, ",")
			for i := range tags {
				switch tags[i] {
				case "update":
					updateFields = append(updateFields, tag)
				}
			}
			if allFields == nil {
				allFields = make([]string, 0, t.NumField())
			}
			allFields = append(allFields, tag)
		}

		entityUpdateFields[name] = updateFields

		entityFields[name] = allFields
	}
}

func GetFields(tableName string, args ...string) ([]string, error) {
	v, ok := entityFields[tableName]
	if !ok {
		return nil, fmt.Errorf("entity with specified table name %s not exist", tableName)
	}

	b := make([]string, len(v))
	copy(b, v)

	return b, nil
}

func GetUpdateFields(tableName string) ([]string, error) {
	v, ok := entityUpdateFields[tableName]
	if !ok {
		return nil, fmt.Errorf("entity with specified table name %s not exist", tableName)
	}

	b := make([]string, len(v))
	copy(b, v)

	return b, nil
}

func GetFieldsPointers(u interface{}, args ...string) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		if len(args) != 0 {
			if val.Type().Field(i).Tag.Get("ops") != args[0] {
				continue
			}
		}
		valueField := val.Field(i)
		v = append(v, valueField.Addr().Interface())
	}
	return v
}
