package infoblog

import (
	"fmt"
	"reflect"
	"strings"
)

var entityFields = map[string][]string{}
var entityUpdateFields = map[string][]string{}
var entityCreateFields = map[string][]string{}
var entityDeleteFields = map[string][]string{}

// register entities
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
		var createFields []string
		var deleteFields []string

		for i := 0; i < t.NumField(); i++ {

			if allFields == nil {
				allFields = make([]string, 0, t.NumField())
			}

			if updateFields == nil {
				updateFields = make([]string, 0, t.NumField())
			}

			if createFields == nil {
				createFields = make([]string, 0, t.NumField())
			}

			if deleteFields == nil {
				deleteFields = make([]string, 0, t.NumField())
			}

			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			field := t.Field(i)

			// Get the field tag value
			filedName := field.Tag.Get("db")
			if filedName == "" || filedName == "-" {
				continue
			}
			allFields = append(allFields, filedName)

			tagsString := field.Tag.Get("ops")
			tags := strings.Split(tagsString, ",")
			if tagsString != "" {
				for i := range tags {
					switch tags[i] {
					case "update":
						updateFields = append(updateFields, filedName)
					case "create":
						createFields = append(createFields, filedName)
					case "delete":
						deleteFields = append(deleteFields, filedName)
					}
				}
			}
		}

		entityUpdateFields[name] = updateFields

		entityFields[name] = allFields

		entityCreateFields[name] = createFields

		entityDeleteFields[name] = deleteFields
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

func GetCreateFields(tableName string) ([]string, error) {
	v, ok := entityCreateFields[tableName]
	if !ok {
		return nil, fmt.Errorf("entity with specified table name %s not exist", tableName)
	}

	b := make([]string, len(v))
	copy(b, v)

	return b, nil
}

func GetDeleteFields(tableName string) ([]string, error) {
	v, ok := entityDeleteFields[tableName]
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
			tagsRaw := val.Type().Field(i).Tag.Get("ops")
			tags := strings.Split(tagsRaw, ",")
			found := false
			for _, tag := range tags {
				if tag == args[0] {
					found = true
				}
			}
			if !found {
				continue
			}
		}
		valueField := val.Field(i)
		v = append(v, valueField.Addr().Interface())
	}
	return v
}
