package db

import (
	"context"
	"fmt"
	"reflect"

	infoblog "gitlab.com/InfoBlogFriends/server"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

type crud struct {
	db *sqlx.DB
}

func (c *crud) create(ctx context.Context, entity interface{}) error {
	ent, ok := entity.(infoblog.Entity)
	if !ok {
		return fmt.Errorf("wrong entity")
	}
	createFields, err := infoblog.GetFields(ent, "create")
	if err != nil {
		return err
	}
	createFieldsPointers := infoblog.GetFieldsPointers(entity, "create")

	queryRaw := sq.Insert(ent.TableName()).Columns(createFields...).Values(createFieldsPointers...)

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	_, err = c.db.MustExecContext(ctx, query, args...).RowsAffected()

	return err
}

func (c *crud) update(ctx context.Context, entity interface{}) error {
	ent, ok := entity.(infoblog.Entity)
	if !ok {
		return fmt.Errorf("wrong entity")
	}
	updateFields, err := infoblog.GetFields(ent, "update")
	if err != nil {
		return err
	}
	updateFieldsPointers := infoblog.GetFieldsPointers(entity, "update")

	val := reflect.ValueOf(entity).Elem()
	uuid := val.FieldByName("UUID").Interface().([]byte)
	queryRaw := sq.Update(ent.TableName()).Where(sq.Eq{"uuid": uuid})
	for i := range updateFields {
		queryRaw = queryRaw.Set(updateFields[i], updateFieldsPointers[i])
	}

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}

func (c *crud) find(ctx context.Context, entity interface{}, dest interface{}) error {
	ent, ok := entity.(infoblog.Entity)
	if !ok {
		return fmt.Errorf("wrong entity")
	}
	fields, err := infoblog.GetFields(ent)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(entity).Elem()
	uuid := val.FieldByName("UUID").Interface().([]byte)
	queryRaw := sq.Select(fields...).From(ent.TableName()).Where(sq.Eq{"uuid": uuid})
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	err = c.db.QueryRowxContext(ctx, query, args...).StructScan(dest)

	return err
}

func (c *crud) list(ctx context.Context, dest interface{}, entity interface{}, limit, offset uint64) error {
	ent, ok := entity.(infoblog.Entity)
	if !ok {
		return fmt.Errorf("wrong entity")
	}
	fields, err := infoblog.GetFields(ent)
	if err != nil {
		return err
	}

	query, args, err := sq.Select(fields...).From(ent.TableName()).Limit(limit).Offset(offset).ToSql()
	if err != nil {
		return err
	}

	if err = c.db.SelectContext(ctx, dest, query, args...); err != nil {
		return err
	}

	return nil
}

//func (c *crud) increment(ctx context.Context, entity interface{}, field string) error {
//	ent, ok := entity.(infoblog.Entity)
//	if !ok {
//		return fmt.Errorf("wrong entity")
//	}
//	fields, err := infoblog.GetFields(ent)
//	if err != nil {
//		return err
//	}
//
//	val := reflect.ValueOf(entity).Elem()
//	uuid := val.FieldByName("UUID").Interface().([]byte)
//
//	queryRaw := sq.Select(fields...).From(ent.TableName()).Where(sq.Eq{"uuid": uuid})
//	query, args, err := queryRaw.ToSql()
//	if err != nil {
//		return err
//	}
//	query = strings.Join([]string{query, "FOR UPDATE"}, " ")
//	tx, err := c.db.Beginx()
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback()
//		} else {
//			_ = tx.Commit()
//		}
//	}()
//
//	if err != nil {
//		return err
//	}
//
//	err = tx.QueryRowxContext(ctx, query, args...).StructScan(entity)
//	if err != nil {
//		return err
//	}
//
//	count := val.FieldByName(field).Interface().(uint64)
//
//	reflect.ValueOf(entity).Elem().FieldByName(field).SetUint(count + 1)
//
//	return err
//}
