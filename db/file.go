package db

import (
	"context"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/types"
)

const (
	activeFile = "UPDATE files SET active = ? WHERE uuid = ?"
)

type filesRepository struct {
	db *sqlx.DB
	crud
}

func NewFilesRepository(db *sqlx.DB) infoblog.FileRepository {
	return &filesRepository{db: db}
}

func (f *filesRepository) Create(ctx context.Context, p *infoblog.File) (int64, error) {
	createFields, err := infoblog.GetCreateFields("files")
	if err != nil {
		return 0, err
	}
	createFieldsPointers := infoblog.GetFieldsPointers(p, "create")

	queryRaw := sq.Insert("files").Columns(createFields...).Values(createFieldsPointers...)
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return 0, err
	}

	return f.db.MustExecContext(ctx, query, args...).RowsAffected()
}

func (f *filesRepository) Update(ctx context.Context, file infoblog.File) error {
	if !file.UUID.Valid {
		return errors.New("wrong file uuid on update")
	}

	updateFields, err := infoblog.GetUpdateFields("files")
	if err != nil {
		return err
	}
	updateFieldsPointers := infoblog.GetFieldsPointers(&file, "update")

	queryRaw := sq.Update("files").Where(sq.Eq{"uuid": file.UUID})
	for i := range updateFields {
		queryRaw = queryRaw.Set(updateFields[i], updateFieldsPointers[i])
	}

	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	res, err := f.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()

	return err
}

func (f *filesRepository) UpdatePostUUID(ctx context.Context, ids []string, post infoblog.Post) error {
	query, args, err := sq.Update("files").Set("foreign_uuid", post.UUID).ToSql()
	if err != nil {
		return err
	}

	queryIn := "SELECT * FROM files WHERE uuid IN (?)"

	uuids := make([]types.NullUUID, 0, len(ids))
	for i := range ids {
		uuids = append(uuids, types.NewNullUUID(ids[i]))
	}
	queryIn, args2, err := sqlx.In(queryIn, uuids)
	if err != nil {
		return err
	}
	s := strings.Split(queryIn, "WHERE")

	query = strings.Join([]string{query, " WHERE", s[1]}, "")

	args = append(args, args2...)

	return f.db.QueryRowContext(ctx, query, args...).Err()
}

func (f *filesRepository) Delete(ctx context.Context, p infoblog.File) error {
	if !p.UUID.Valid {
		return errors.New("repository delete wrong file id")
	}
	return f.db.QueryRowContext(ctx, activeFile, 0, p.UUID).Err()
}

func (f *filesRepository) Find(ctx context.Context, file infoblog.File) (infoblog.File, error) {
	fields, err := infoblog.GetFields(&infoblog.File{})
	if err != nil {
		return infoblog.File{}, err
	}

	queryRaw := sq.Select(fields...).From("files").Where(sq.Eq{"uuid": file.UUID})
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return infoblog.File{}, err
	}

	file = infoblog.File{}
	err = f.db.QueryRowxContext(ctx, query, args...).StructScan(&file)

	return file, err
}

func (f *filesRepository) FindAll(ctx context.Context, postUUID string) ([]infoblog.File, error) {
	return f.FindByTypeFUUID(ctx, 1, postUUID)
}

func (f *filesRepository) FindByTypeFUUID(ctx context.Context, typeID int64, foreignUUID string) ([]infoblog.File, error) {

	fields, err := infoblog.GetFields(&infoblog.File{})
	if err != nil {
		return nil, err
	}

	queryRaw := sq.Select(fields...).From("files").Where(sq.Eq{"type": typeID, "foreign_uuid": foreignUUID})
	query, args, err := queryRaw.ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := f.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := make([]infoblog.File, 0)

	for rows.Next() {
		file := infoblog.File{}
		err = rows.StructScan(&file)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (f filesRepository) FindByPostsIDs(ctx context.Context, postsIDs []string) ([]infoblog.File, error) {
	uuids := make([]types.NullUUID, 0, len(postsIDs))
	for i := range postsIDs {
		uuids = append(uuids, types.NewNullUUID(postsIDs[i]))
	}

	fields, err := infoblog.GetFields(&infoblog.File{})
	if err != nil {
		return nil, err
	}

	queryRaw := sq.Select(fields...).From("files").Where(sq.Eq{"active": 1, "type": 1})
	var args []interface{}
	query, _, err := queryRaw.ToSql()
	if err != nil {
		return nil, err
	}
	query = query + " AND foreign_uuid IN (?)"
	query, args, err = sqlx.In(query, 1, 1, uuids)

	if err != nil {
		return nil, err
	}
	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = f.db.Rebind(query)

	rows, err := f.db.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]infoblog.File, 0)

	for rows.Next() {
		file := infoblog.File{}
		err = rows.StructScan(&file)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (f *filesRepository) FindByIDs(ctx context.Context, ids []string) ([]infoblog.File, error) {
	fields, err := infoblog.GetFields(&infoblog.File{})
	if err != nil {
		return nil, err
	}

	query, _, err := sq.Select(fields...).From("files").ToSql()
	if err != nil {
		return nil, err
	}
	query = strings.Join([]string{query, " WHERE uuid IN (?)"}, "")

	query, args, err := sqlx.In(query, ids)
	if err != nil {
		return nil, err
	}
	var files []infoblog.File

	err = f.db.Select(&files, query, args...)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (f *filesRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.File, error) {
	var events []infoblog.File
	err := f.crud.listx(ctx, &events, infoblog.File{}, condition)
	if err != nil {
		return nil, err
	}

	return events, nil
}
