package db

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	createFile = "INSERT INTO files (type, foreign_id, dir, name, user_id, user_uuid, uuid, foreign_uuid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	updateFile = "UPDATE files SET type = ?, foreign_id = ?, dir = ?, name = ?, user_id = ? WHERE id = ?"
	activeFile = "UPDATE files SET active = ? WHERE id = ?"

	selectFile     = "SELECT id, type, foreign_id, dir, name, active, user_id, created_at, updated_at FROM files WHERE active = 1 AND id = ?"
	findAllFile    = "SELECT id, type, foreign_id, dir, name, active, user_id, created_at, updated_at FROM files WHERE active = 1 AND type = ? AND foreign_id = ?"
	findByPostsIDs = "SELECT id, type, foreign_id, dir, name, active, user_id, uuid, created_at, updated_at FROM files WHERE active = 1 AND type = 1 AND foreign_id IN (?)"
)

type filesRepository struct {
	db *sqlx.DB
}

func NewFilesRepository(db *sqlx.DB) infoblog.FileRepository {
	return &filesRepository{db: db}
}

func (f *filesRepository) Create(ctx context.Context, p *infoblog.File) (int64, error) {
	res, err := f.db.ExecContext(ctx, createFile, p.Type, p.ForeignID, p.Dir, p.Name, p.UserID, p.UserUUID, p.UUID, p.ForeignUUID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (f *filesRepository) Update(ctx context.Context, p infoblog.File) error {
	if p.ID < 1 {
		return errors.New("repository update wrong file id")
	}
	return f.db.QueryRowContext(ctx, updateFile, p.Type, p.ForeignID, p.Dir, p.Name, p.UserID).Err()
}

func (f *filesRepository) Delete(ctx context.Context, p infoblog.File) error {
	if p.ID < 1 {
		return errors.New("repository delete wrong file id")
	}
	return f.db.QueryRowContext(ctx, activeFile, 0, p.ID).Err()
}

func (f *filesRepository) Find(ctx context.Context, id int64) (infoblog.File, error) {
	if id < 1 {
		return infoblog.File{}, errors.New("repository find wrong file id")
	}
	file := infoblog.File{}
	err := f.db.QueryRowContext(ctx, selectFile, id).Scan(&file)

	return file, err
}

func (f *filesRepository) FindAll(ctx context.Context, postID int64) ([]infoblog.File, error) {
	return f.FindByTypeFID(ctx, 1, postID)
}

func (f *filesRepository) FindByTypeFID(ctx context.Context, typeID int64, foreignID int64) ([]infoblog.File, error) {
	rows, err := f.db.QueryContext(ctx, findAllFile, typeID, foreignID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	files := make([]infoblog.File, 0)

	for rows.Next() {
		file := infoblog.File{}
		err = rows.Scan(&file.ID, file.Type, file.ForeignID, file.Dir, file.Name, file.Active, file.UserID, file.CreatedAt, file.UpdatedAt)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

func (f filesRepository) FindByPostsIDs(ctx context.Context, postsIDs []int) ([]infoblog.File, error) {
	if len(postsIDs) < 1 {
		return nil, nil
	}
	query, args, err := sqlx.In(findByPostsIDs, postsIDs)

	if err != nil {
		return nil, err
	}
	// sqlx.In returns queries with the `?` bindvar, we can rebind it for our backend
	query = f.db.Rebind(query)
	rows, err := f.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := make([]infoblog.File, 0)

	for rows.Next() {
		file := infoblog.File{}
		err = rows.Scan(&file.ID, &file.Type, &file.ForeignID, &file.Dir, &file.Name, &file.Active, &file.UserID, &file.UUID, &file.CreatedAt, &file.UpdatedAt)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
