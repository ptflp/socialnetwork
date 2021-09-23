package migration

import (
	"strings"

	"gitlab.com/InfoBlogFriends/server/config"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

//go:generate qtc -dir=./

const selectTableFields = "SELECT COLUMN_NAME  FROM INFORMATION_SCHEMA.COLUMNS  WHERE  TABLE_SCHEMA = ? AND TABLE_NAME = ?"

type Migrator struct {
	db     *sqlx.DB
	dbConf config.DB
}

func NewMigrator(db *sqlx.DB, dbConf config.DB) *Migrator {
	return &Migrator{db: db, dbConf: dbConf}
}

func (m *Migrator) Migrate() error {
	tables := infoblog.GetTables()
	var err error
	for name := range tables {
		table := tables[name]
		var tableFields []string
		err = m.db.Select(&tableFields, selectTableFields, m.dbConf.DBName, table.Name)
		if err != nil {
			return err
		}
		tableFieldsMap := make(map[string]string, len(tableFields))
		for i := range tableFields {
			tableFieldsMap[tableFields[i]] = tableFields[i]
		}
		if len(tableFields) < 1 {
			createQuery := CreateTable(table)
			queries := strings.Split(createQuery, ";")
			for i := range queries {
				queries[i] = strings.TrimSpace(queries[i])
				if queries[i] == "" {
					continue
				}
				_, err = m.db.Queryx(queries[i])
				if err != nil {
					return err
				}
			}
		}
		if len(tableFields) > 0 {
			entityFields, _ := infoblog.GetFields(table.Entity)
			var diff map[string]infoblog.Field
			for i := range entityFields {
				if _, ok := tableFieldsMap[entityFields[i]]; !ok {
					if diff == nil {
						diff = make(map[string]infoblog.Field, len(entityFields))
					}
					diff[entityFields[i]] = table.FieldsMap[entityFields[i]]
				}
			}
			for fieldName := range diff {
				if fieldName == "" {
					continue
				}
				alterQuery := AlterTable(diff[fieldName])
				queries := strings.Split(alterQuery, ";")
				for i := range queries {
					queries[i] = strings.TrimSpace(queries[i])
					if queries[i] == "" {
						continue
					}
					_, err = m.db.Queryx(queries[i])
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return err
}
