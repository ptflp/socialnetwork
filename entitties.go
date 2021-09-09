package infoblog

import (
	"fmt"
	"reflect"
	"strings"

	sq "github.com/Masterminds/squirrel"
)

//go:generate go get -u github.com/valyala/quicktemplate/qtc
//go:generate qtc -dir=migration

var (
	entityFields       map[string][]string
	entityUpdateFields map[string][]string
	entityCreateFields map[string][]string
	entityCountFields  map[string][]string
	tables             map[string]Table
)

type Table struct {
	Name        string
	Fields      []Field
	FieldsMap   map[string]Field
	Constraints []Constraint
	Entity      Tabler
}

func (t Table) CreateQuery() string {
	return ""
}

type Field struct {
	Name       string
	Type       string
	Default    string
	Constraint Constraint
	TableName  string
}

type Constraint struct {
	Index     bool
	Unique    bool
	FieldName string
}

type Constraints struct {
	Name string
	Type string
}

// register entities
func init() {
	RegisterEntities(
		User{},
		Like{},
		HashTag{},
		Subscriber{},
		File{},
		PostEntity{},
		Chat{},
		ChatMessages{},
		ChatParticipant{},
		ChatPrivateUser{},
		Comment{},
		Friend{},
		Moderate{},
		HashTag{},
	)
}

// register entities
func RegisterEntities(entities ...Tabler) {
	tableEntities := make(map[string]Tabler, len(entities))
	for i := range entities {
		tableEntities[entities[i].TableName()] = entities[i]
	}

	tables = make(map[string]Table, len(entities))
	entityFields = make(map[string][]string, len(entities))
	entityUpdateFields = make(map[string][]string, len(entities))
	entityCreateFields = make(map[string][]string, len(entities))
	entityCountFields = make(map[string][]string, len(entities))

	for name, entity := range tableEntities {
		table := Table{
			Name:   name,
			Entity: entity,
		}
		t := reflect.TypeOf(entity)

		var allFields []string
		var updateFields []string
		var createFields []string
		var countFields []string

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

			if countFields == nil {
				countFields = make([]string, 0, t.NumField())
			}

			if table.FieldsMap == nil {
				table.FieldsMap = make(map[string]Field, t.NumField())
			}

			// Get the field, returns https://golang.org/pkg/reflect/#StructField
			structField := t.Field(i)
			// Get the structField tag value
			fieldName := structField.Tag.Get("db")

			if fieldName == "" || fieldName == "-" {
				continue
			}
			allFields = append(allFields, fieldName)

			field := Field{
				Name:      fieldName,
				Type:      structField.Tag.Get("orm_type"),
				Default:   structField.Tag.Get("orm_default"),
				TableName: table.Name,
			}
			constraintRaw := structField.Tag.Get("orm_index")
			constraintPieces := strings.Split(constraintRaw, ",")
			if len(constraintPieces) < 1 {
				field.Constraint = Constraint{}
			}
			if len(constraintPieces) > 0 {
				for i := range constraintPieces {
					switch constraintPieces[i] {
					case "index":
						field.Constraint.Index = true
					case "unique":
						field.Constraint.Unique = true
					}
				}
			}
			if field.Constraint.Index {
				field.Constraint.FieldName = field.Name
				table.Constraints = append(table.Constraints, field.Constraint)
			}
			table.Fields = append(table.Fields, field)
			table.FieldsMap[field.Name] = field

			opsRaw := structField.Tag.Get("ops")
			ops := strings.Split(opsRaw, ",")
			if opsRaw != "" {
				for i := range ops {
					switch ops[i] {
					case "update":
						updateFields = append(updateFields, fieldName)
					case "create":
						createFields = append(createFields, fieldName)
					case "count":
						countFields = append(countFields, fieldName)
					}
				}
			}
		}

		entityUpdateFields[name] = updateFields

		entityFields[name] = allFields

		entityCreateFields[name] = createFields

		entityCountFields[name] = countFields

		tables[name] = table
	}
}

func GetFields(entity Tabler, args ...string) ([]string, error) {
	if len(args) < 1 {
		return GetAllFields(entity.TableName())
	}
	switch args[0] {
	case "create":
		return GetCreateFields(entity.TableName())
	case "update":
		return GetUpdateFields(entity.TableName())
	case "count":
		return GetCountFields(entity.TableName())
	default:
		return GetAllFields(entity.TableName())
	}
}

func GetAllFields(tableName string, args ...string) ([]string, error) {
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

func GetCountFields(tableName string) ([]string, error) {
	v, ok := entityCountFields[tableName]
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

func GetTables() map[string]Table {
	return tables
}

type In struct {
	Field string
	Args  []interface{}
}

type Other struct {
	Condition string
	Args      []interface{}
}

type Order struct {
	Field string
	Asc   bool
}

type LimitOffset struct {
	Offset int64
	Limit  int64
}

type Condition struct {
	Equal       *sq.Eq
	In          *In
	NotIn       *In
	Other       *Other
	Order       *Order
	LimitOffset *LimitOffset
}
