package db

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"

	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

type crud struct {
	db *sqlx.DB
}

func (c *crud) create(ctx context.Context, entity infoblog.Tabler) error {
	createFields, err := infoblog.GetFields(entity, "create")
	if err != nil {
		return err
	}
	createFieldsPointers := infoblog.GetFieldsPointers(entity, "create")

	queryRaw := sq.Insert(entity.TableName()).Columns(createFields...).Values(createFieldsPointers...)

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx, query, args...)

	return err
}

func (c *crud) update(ctx context.Context, entity infoblog.Tabler) error {
	ent := entity
	updateFields, err := infoblog.GetFields(ent, "update")
	if err != nil {
		return err
	}
	updateFieldsPointers := infoblog.GetFieldsPointers(entity, "update")

	val := reflect.ValueOf(entity).Elem()
	uuid := val.FieldByName("UUID").Interface().(types.NullUUID)
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

func (c *crud) find(ctx context.Context, entity infoblog.Tabler, dest interface{}) error {
	ent := entity
	fields, err := infoblog.GetFields(ent)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(entity).Elem()
	uuid := val.FieldByName("UUID").Interface().(types.NullUUID)
	queryRaw := sq.Select(fields...).From(ent.TableName()).Where(sq.Eq{"uuid": uuid})
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	err = c.db.QueryRowxContext(ctx, query, args...).StructScan(dest)

	return err
}

func (c *crud) first(ctx context.Context, dest interface{}) error {
	ent, ok := dest.(infoblog.Tabler)
	if !ok {
		return fmt.Errorf("wrong entity")
	}
	fields, err := infoblog.GetFields(ent)
	if err != nil {
		return err
	}

	queryRaw := sq.Select(fields...).From(ent.TableName()).Limit(1)
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}

	err = c.db.QueryRowxContext(ctx, query, args...).StructScan(dest)

	return err
}

func (c *crud) getCount(ctx context.Context, entity infoblog.Tabler, condition infoblog.Condition) (uint64, error) {
	ent := entity
	var whereState bool

	queryRaw := sq.Select("COUNT(*)").From(ent.TableName())

	if condition.Equal != nil {
		queryRaw = queryRaw.Where(condition.Equal)
		whereState = true
	}

	if condition.Other != nil {
		queryRaw = queryRaw.Where(condition.Other.Condition, condition.Other.Args...)
		whereState = true
	}

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return 0, err
	}

	if condition.In != nil {
		queryIn := fmt.Sprintf("SELECT * FROM files WHERE %s IN (?)", condition.In.Field)

		queryIn, inArgs, err := sqlx.In(queryIn, condition.In.Args)
		if err != nil {
			return 0, err
		}

		s := strings.Split(queryIn, "WHERE")
		sep := " WHERE"

		if whereState {
			sep = " AND"
		}

		query = strings.Join([]string{query, sep, s[1]}, "")

		args = append(args, inArgs...)
		whereState = true
	}

	if condition.NotIn != nil {
		queryNotIn := fmt.Sprintf("SELECT * FROM files WHERE %s IN (?)", condition.NotIn.Field)

		queryNotIn, inArgs, err := sqlx.In(queryNotIn, condition.NotIn.Args)
		if err != nil {
			return 0, err
		}

		queryNotIn = strings.Replace(queryNotIn, "IN (", "NOT IN (", 1)

		s := strings.Split(queryNotIn, "WHERE")
		sep := " WHERE"

		if whereState {
			sep = " AND"
		}

		query = strings.Join([]string{query, sep, s[1]}, "")

		args = append(args, inArgs...)
	}

	query = strings.Join([]string{query, " LIMIT 1"}, "")

	rows, err := c.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	var count uint64
	// iterate over each row
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	// check the error from rows
	err = rows.Err()

	return count, err
}

func (c *crud) list(ctx context.Context, dest interface{}, entity infoblog.Tabler, limit, offset uint64) error {
	ent := entity
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

func (c *crud) listx(ctx context.Context, dest interface{}, entity infoblog.Tabler, condition infoblog.Condition) error {
	fields, err := infoblog.GetFields(entity)
	if err != nil {
		return err
	}
	var whereState bool

	sQuery := sq.Select(fields...).From(entity.TableName())

	if condition.Equal != nil {
		sQuery = sQuery.Where(condition.Equal)
		whereState = true
	}

	if condition.Other != nil {
		sQuery = sQuery.Where(condition.Other.Condition, condition.Other.Args...)
		whereState = true
	}

	query, args, err := sQuery.ToSql()
	if err != nil {
		return err
	}

	if condition.In != nil {
		queryIn := fmt.Sprintf("SELECT * FROM files WHERE %s IN (?)", condition.In.Field)

		queryIn, inArgs, err := sqlx.In(queryIn, condition.In.Args)
		if err != nil {
			return err
		}

		s := strings.Split(queryIn, "WHERE")
		sep := " WHERE"

		if whereState {
			sep = " AND"
		}

		query = strings.Join([]string{query, sep, s[1]}, "")

		args = append(args, inArgs...)
		whereState = true
	}

	if condition.NotIn != nil {
		queryNotIn := fmt.Sprintf("SELECT * FROM files WHERE %s IN (?)", condition.NotIn.Field)

		queryNotIn, inArgs, err := sqlx.In(queryNotIn, condition.NotIn.Args)
		if err != nil {
			return err
		}

		queryNotIn = strings.Replace(queryNotIn, "IN (", "NOT IN (", 1)

		s := strings.Split(queryNotIn, "WHERE")
		sep := " WHERE"

		if whereState {
			sep = " AND"
		}

		query = strings.Join([]string{query, sep, s[1]}, "")

		args = append(args, inArgs...)
	}

	if condition.Order != nil {
		direction := "DESC"
		if condition.Order.Asc {
			direction = "ASC"
		}
		order := fmt.Sprintf(" ORDER BY %s %s", condition.Order.Field, direction)
		query = strings.Join([]string{query, order}, "")
	}

	if condition.LimitOffset != nil {
		limitOffset := fmt.Sprintf(" LIMIT %d OFFSET %d", condition.LimitOffset.Limit, condition.LimitOffset.Offset)

		query = strings.Join([]string{query, limitOffset}, "")
	}

	if condition.ForUpdate {
		query = query + " FOR UPDATE"
	}

	if err = c.db.SelectContext(ctx, dest, query, args...); err != nil {
		return err
	}

	return nil
}

func (c *crud) updatex(ctx context.Context, entity infoblog.Tabler, condition infoblog.Condition) error {
	ent := entity
	updateFields, err := infoblog.GetFields(ent, "update")
	if err != nil {
		return err
	}
	var whereState bool

	updateFieldsPointers := infoblog.GetFieldsPointers(entity, "update")

	updateRaw := sq.Update(ent.TableName())

	if condition.Equal != nil {
		updateRaw = updateRaw.Where(condition.Equal)
		whereState = true
	}

	if condition.Other != nil {
		updateRaw = updateRaw.Where(condition.Other.Condition, condition.Other.Args...)
		whereState = true
	}

	for i := range updateFields {
		updateRaw = updateRaw.Set(updateFields[i], updateFieldsPointers[i])
	}

	query, args, err := updateRaw.ToSql()
	if err != nil {
		return err
	}

	if condition.In != nil {
		queryIn := fmt.Sprintf("SELECT * FROM files WHERE %s IN (?)", condition.In.Field)

		queryIn, inArgs, err := sqlx.In(queryIn, condition.In.Args)
		if err != nil {
			return err
		}

		s := strings.Split(queryIn, "WHERE")
		sep := " WHERE"

		if whereState {
			sep = " AND"
		}

		query = strings.Join([]string{query, sep, s[1]}, "")

		args = append(args, inArgs...)
		whereState = true
	}

	if condition.NotIn != nil {
		queryNotIn := fmt.Sprintf("SELECT * FROM files WHERE %s IN (?)", condition.NotIn.Field)

		queryNotIn, inArgs, err := sqlx.In(queryNotIn, condition.NotIn.Args)
		if err != nil {
			return err
		}

		queryNotIn = strings.Replace(queryNotIn, "IN (", "NOT IN (", 1)

		s := strings.Split(queryNotIn, "WHERE")
		sep := " WHERE"

		if whereState {
			sep = " AND"
		}

		query = strings.Join([]string{query, sep, s[1]}, "")

		args = append(args, inArgs...)
	}

	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}

func (c *crud) count(ctx context.Context, entity infoblog.Tabler, field, ops string) error {
	switch ops {
	case "decr":
	case "incr":
		break
	default:
		return fmt.Errorf("bad count operation")
	}

	ent := entity
	fields, err := infoblog.GetFields(ent)
	if err != nil {
		return err
	}

	val := reflect.ValueOf(entity).Elem()
	uuid := val.FieldByName("UUID").Interface().(types.NullUUID)

	queryRaw := sq.Select(fields...).From(ent.TableName()).Where(sq.Eq{"uuid": uuid})
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	query = strings.Join([]string{query, "FOR UPDATE"}, " ")
	tx, err := c.db.Beginx()
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	if err != nil {
		return err
	}

	err = tx.QueryRowxContext(ctx, query, args...).StructScan(entity)
	if err != nil {
		return err
	}
	field = strings.Title(field)

	count := val.FieldByName(field).Interface().(types.NullUint64)

	switch ops {
	case "decr":
		if count.Uint64.Uint64 < 1 {
			count = types.NullUint64{}
		} else {
			count.Uint64.Uint64--
			count.Valid = true
		}
	case "incr":
		count.Uint64.Uint64++
		count.Valid = true
	default:
		return fmt.Errorf("bad count operation %s", ops)
	}

	v := reflect.ValueOf(count)

	reflect.ValueOf(entity).Elem().FieldByName(field).Set(v)

	updateFields, err := infoblog.GetFields(ent, "count")
	if err != nil {
		return err
	}
	updateFieldsPointers := infoblog.GetFieldsPointers(entity, "count")

	updateRaw := sq.Update(ent.TableName()).Where(sq.Eq{"uuid": uuid})
	for i := range updateFields {
		updateRaw = updateRaw.Set(updateFields[i], updateFieldsPointers[i])
	}

	query, args, err = updateRaw.ToSql()
	if err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}
