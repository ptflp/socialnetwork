package infoblog

import (
	"fmt"
	"reflect"
)

var entityFields = map[string][]string{}

func init() {

	entities := map[string]interface{}{
		"users":      User{},
		"likes":      Like{},
		"hashtags":   HashTag{},
		"subscribes": Subscriber{},
		"files":      File{},
	}

	for name, entity := range entities {
		t := reflect.TypeOf(entity)

		var fields []string

		for i := 0; i < t.NumField(); i++ {
			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			field := t.Field(i)

			// Get the field tag value
			tag := field.Tag.Get("db")
			if fields == nil {
				fields = make([]string, 0, t.NumField())
			}
			fields = append(fields, tag)
		}

		entityFields[name] = fields
	}
}

func GetFields(tableName string) ([]string, error) {
	v, ok := entityFields[tableName]
	if !ok {
		return nil, fmt.Errorf("entity with specified table name %s not exist", tableName)
	}

	return v, nil
}
