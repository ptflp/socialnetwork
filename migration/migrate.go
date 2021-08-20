package migration

import (
	"strings"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

//go:generate qtc -dir=./

const tableFields = "SELECT COLUMN_NAME  FROM INFORMATION_SCHEMA.COLUMNS  WHERE  TABLE_SCHEMA = ? AND TABLE_NAME = ?"

type Migrator struct {
	db *sqlx.DB
}

func NewMigrator(db *sqlx.DB) *Migrator {
	return &Migrator{db: db}
}

func (m *Migrator) Migrate() error {
	tables := infoblog.GetTables()
	var err error
	for name := range tables {
		table := tables[name]
		var fields []string
		err = m.db.Select(&fields, tableFields, "infoblog", table.Name)
		if err != nil {
			return err
		}
		if len(fields) < 1 {
			createQuery := CreateTable(table)
			queries := strings.Split(createQuery, ";")
			for i := range queries {
				if queries[i] == "" {
					continue
				}
				queries[i] = strings.TrimSpace(queries[i])
				_, err = m.db.Queryx(queries[i])
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}
