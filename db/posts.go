package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	createPost     = "INSERT INTO posts (body, user_uuid, uuid, type, price, active) VALUES (?, ?, ?, ?, ?, ?)"
	updatePost     = "UPDATE posts SET body = ?, active = ?, price = ? WHERE uuid = ?"
	deletePost     = "UPDATE posts SET active = ? WHERE uuid = ?"
	countAllRecent = "SELECT COUNT(p.uuid) FROM posts p WHERE p.active = 1"
)

type postsRepository struct {
	db *sqlx.DB
	crud
}

func (pr *postsRepository) Create(ctx context.Context, p infoblog.Post) (int64, error) {
	res, err := pr.db.ExecContext(ctx, createPost, p.Body, p.UserUUID, p.UUID, p.Type, p.Price, p.Active)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (pr *postsRepository) Update(ctx context.Context, p infoblog.Post) error {
	if !p.UUID.Valid {
		return errors.New("repository wrong post id")
	}
	_, err := pr.db.MustExecContext(ctx, updatePost, p.Body, p.Active, p.Price, p.UUID).RowsAffected()

	return err
}

func (pr *postsRepository) Delete(ctx context.Context, p infoblog.Post) error {
	if !p.UUID.Valid {
		return errors.New("repository wrong post id")
	}
	_, err := pr.db.MustExecContext(ctx, deletePost, 0, p.UUID).RowsAffected()

	return err
}

func (pr *postsRepository) Find(ctx context.Context, p infoblog.Post) (infoblog.Post, error) {
	if !p.UUID.Valid {
		return infoblog.Post{}, errors.New("repository wrong post uuid")
	}

	fields, err := infoblog.GetFields(&p)
	if err != nil {
		return infoblog.Post{}, err
	}

	query, args, err := sq.Select(fields...).From("posts").Where(sq.Eq{"uuid": p.PostEntity.UUID}).ToSql()
	if err != nil {
		return infoblog.Post{}, err
	}

	if err := pr.db.QueryRowxContext(ctx, query, args...).StructScan(&p.PostEntity); err != nil {
		return infoblog.Post{}, err
	}

	return infoblog.Post{
		PostEntity: p.PostEntity,
	}, nil
}

func (pr *postsRepository) FindAll(ctx context.Context, user infoblog.User, limit int64, offset int64) ([]infoblog.Post, map[string]int, []string, error) {
	fields, err := infoblog.GetFields(&infoblog.PostEntity{})
	if err != nil {
		return nil, nil, nil, err
	}

	for i := range fields {
		s := strings.Join([]string{"p", fields[i]}, ".")
		fields[i] = s
	}

	userFields, err := infoblog.GetFields(&infoblog.User{})
	if err != nil {
		return nil, nil, nil, err
	}

	for i := range userFields {
		s := strings.Join([]string{"u", userFields[i]}, ".")
		userFields[i] = s
	}
	fields = append(fields, userFields...)

	query, args, err := sq.Select(fields...).From("posts p").LeftJoin("users u on p.user_uuid = u.uuid").Where(sq.Eq{"p.active": 1, "p.user_uuid": user.UUID}).OrderBy("p.created_at DESC").Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := pr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	postDataRes := make([]infoblog.Post, 0, limit)
	postIdIndexMap := make(map[string]int)
	postsIDs := make([]string, 0, limit)

	for rows.Next() {
		post := infoblog.Post{}
		pFieldsPointers := infoblog.GetFieldsPointers(&post.PostEntity)
		uFieldsPointers := infoblog.GetFieldsPointers(&post.User)

		pFieldsPointers = append(pFieldsPointers, uFieldsPointers...)
		err = rows.Scan(pFieldsPointers...)
		if err != nil {
			return nil, nil, nil, err
		}

		postsIDs = append(postsIDs, post.UUID.String)
		postDataRes = append(postDataRes, post)
		postIdIndexMap[post.UUID.String] = len(postDataRes) - 1
	}

	return postDataRes, postIdIndexMap, postsIDs, nil
}

func (pr *postsRepository) FindAllRecent(ctx context.Context, limit, offset int64) ([]infoblog.Post, map[string]int, []string, error) {
	fields, err := infoblog.GetFields(&infoblog.PostEntity{})
	if err != nil {
		return nil, nil, nil, err
	}

	for i := range fields {
		s := strings.Join([]string{"p", fields[i]}, ".")
		fields[i] = s
	}

	userFields, err := infoblog.GetFields(&infoblog.User{})
	if err != nil {
		return nil, nil, nil, err
	}

	for i := range userFields {
		s := strings.Join([]string{"u", userFields[i]}, ".")
		userFields[i] = s
	}
	fields = append(fields, userFields...)

	query, args, err := sq.Select(fields...).From("posts p").LeftJoin("users u on p.user_uuid = u.uuid").Where(sq.Eq{"p.Active": 1}).OrderBy("p.created_at DESC").Limit(uint64(limit)).Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, nil, nil, err
	}

	rows, err := pr.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	postDataRes := make([]infoblog.Post, 0, limit)
	postUUIDIndexMap := make(map[string]int)
	postsUUID := make([]string, 0, limit)

	for rows.Next() {
		post := infoblog.Post{}
		pFieldsPointers := infoblog.GetFieldsPointers(&post.PostEntity)
		uFieldsPointers := infoblog.GetFieldsPointers(&post.User)

		pFieldsPointers = append(pFieldsPointers, uFieldsPointers...)
		err = rows.Scan(pFieldsPointers...)
		if err != nil {
			return nil, nil, nil, err
		}

		postsUUID = append(postsUUID, post.UUID.String)
		postDataRes = append(postDataRes, post)
		postUUIDIndexMap[post.UUID.String] = len(postDataRes) - 1
	}

	return postDataRes, postUUIDIndexMap, postsUUID, nil
}

func (pr *postsRepository) CountRecent(ctx context.Context) (int64, error) {

	var count sql.NullInt64
	err := pr.db.QueryRowContext(ctx, countAllRecent).Scan(&count)

	return count.Int64, err
}

func (pr *postsRepository) CountByUser(ctx context.Context, user infoblog.User) (int64, error) {

	var count sql.NullInt64

	query, args, err := sq.Select("COUNT(uuid)").From("posts").Where(sq.Eq{"user_uuid": user.UUID, "active": 1}).ToSql()
	if err != nil {
		return count.Int64, err
	}

	err = pr.db.QueryRowContext(ctx, query, args...).Scan(&count)

	return count.Int64, err
}

func (pr *postsRepository) Count(ctx context.Context, p infoblog.Post, field, ops string) (infoblog.Post, error) {
	err := pr.count(ctx, &p.PostEntity, field, ops)
	if err != nil {
		return infoblog.Post{}, err
	}

	return p, nil
}

func (pr *postsRepository) First(ctx context.Context) (infoblog.Post, error) {
	var p infoblog.Post
	err := pr.first(ctx, &p.PostEntity)
	if err != nil {
		return infoblog.Post{}, err
	}

	return p, nil
}

func (pr *postsRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.PostEntity, error) {
	var posts []infoblog.PostEntity
	err := pr.crud.listx(ctx, &posts, infoblog.PostEntity{}, condition)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (pr *postsRepository) GetCount(ctx context.Context, condition infoblog.Condition) (uint64, error) {
	return pr.getCount(ctx, infoblog.PostEntity{}, condition)
}

func NewPostsRepository(db *sqlx.DB) infoblog.PostRepository {
	return &postsRepository{db: db, crud: crud{db: db}}
}
