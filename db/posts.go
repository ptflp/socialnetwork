package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	createPost = "INSERT INTO posts (body, file_id, user_id, uuid, file_uuid, type) VALUES (?, ?, ?, ?, ?, 1)"

	countActivePosts = "SELECT COUNT(id) as count From posts WHERE active = 1"

	updatePost = "UPDATE posts SET body = ?, file_id = ?, active = ? WHERE id = ? AND user_id = ?"
	deletePost = "UPDATE posts SET active = ? WHERE id = ?"

	findPost      = "SELECT type, body, user_id, active, file_id, created_at, updated_at FROM posts WHERE id = ? AND type = 1"
	findAllPost   = "SELECT id, type, body, active, file_id, created_at, updated_at FROM posts WHERE user_id = ? AND active = 1"
	findAllRecent = "SELECT p.id, p.type, body, p.active, p.created_at, p.updated_at, u.id as uid, u.name, u.second_name, u.email, u.phone FROM posts p LEFT JOIN users u on p.user_id = u.id WHERE p.active = 1 AND p.type = 1 ORDER BY p.created_at LIMIT ? OFFSET ?"

	countAllRecent = "SELECT COUNT(p.id) FROM posts p WHERE p.active = 1 AND p.type = 1"
)

type postsRepository struct {
	db *sqlx.DB
}

func (pr *postsRepository) Create(ctx context.Context, p infoblog.Post) (int64, error) {
	res, err := pr.db.ExecContext(ctx, createPost, p.Body, p.FileID, p.UserID, p.UUID, p.FileUUID)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (pr *postsRepository) Update(ctx context.Context, p infoblog.Post) error {
	if p.ID == 0 {
		return errors.New("repository wrong post id")
	}
	_, err := pr.db.MustExecContext(ctx, updatePost, p.Body, p.FileID, p.Active, p.ID, p.UserID).RowsAffected()

	return err
}

func (pr *postsRepository) Delete(ctx context.Context, p infoblog.Post) error {
	if p.ID == 0 {
		return errors.New("repository wrong post id")
	}
	_, err := pr.db.MustExecContext(ctx, deletePost, p.Active, p.ID).RowsAffected()

	return err
}

func (pr *postsRepository) Find(ctx context.Context, id int64) (infoblog.Post, error) {
	if id < 1 {
		return infoblog.Post{}, errors.New("repository wrong post id")
	}

	var (
		typeID    sql.NullInt64
		body      sql.NullString
		fileID    sql.NullInt64
		active    sql.NullInt64
		userID    sql.NullInt64
		createdAt sql.NullTime
		updatedAt sql.NullTime
	)

	if err := pr.db.QueryRowContext(ctx, findPost, id).Scan(&typeID, &body, &userID, &active, &fileID, &createdAt, &updatedAt); err != nil {
		return infoblog.Post{}, err
	}

	return infoblog.Post{
		PostEntity: infoblog.PostEntity{
			ID:        id,
			Type:      typeID.Int64,
			Body:      body.String,
			FileID:    fileID.Int64,
			UserID:    userID.Int64,
			Active:    active.Int64,
			CreatedAt: createdAt.Time,
			UpdatedAt: updatedAt.Time,
		},
	}, nil
}

func (pr *postsRepository) FindAll(ctx context.Context, uid int64) ([]infoblog.Post, map[int64]int, []int, error) {
	if uid < 1 {
		return nil, nil, nil, errors.New("repository wrong post id")
	}

	rows, err := pr.db.QueryContext(ctx, findAllRecent, uid)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	posts := make([]infoblog.Post, 0)

	for rows.Next() {
		post := infoblog.Post{}
		err = rows.Scan(&post.ID, &post.Type, &post.Body, &post.Active, &post.FileID, &post.UserID, &post.CreatedAt, &post.UpdatedAt)
		post.UserID = uid

		if err != nil {
			return nil, nil, nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil, nil, nil
}

func (pr *postsRepository) FindAllRecent(ctx context.Context, limit, offset int64) ([]infoblog.Post, map[int64]int, []int, error) {
	rows, err := pr.db.QueryContext(ctx, findAllRecent, limit, offset)
	if err != nil {
		return nil, nil, nil, err
	}

	defer rows.Close()

	postDataRes := make([]infoblog.Post, 0, limit)
	postIdIndexMap := make(map[int64]int)
	postsIDs := make([]int, 0, limit)

	for rows.Next() {
		post := infoblog.Post{}
		err = rows.Scan(&post.ID, &post.Type, &post.Body, &post.Active, &post.CreatedAt, &post.UpdatedAt, &post.User.ID, &post.User.Name, &post.User.SecondName, &post.User.Email, &post.User.Phone)
		if err != nil {
			return nil, nil, nil, err
		}

		postsIDs = append(postsIDs, int(post.ID))
		postDataRes = append(postDataRes, post)
		postIdIndexMap[post.ID] = len(postDataRes) - 1
	}

	return postDataRes, postIdIndexMap, postsIDs, nil
}

func (pr *postsRepository) CountRecent(ctx context.Context) (int64, error) {

	var count sql.NullInt64
	err := pr.db.QueryRowContext(ctx, countAllRecent).Scan(&count)

	return count.Int64, err
}

func NewPostsRepository(db *sqlx.DB) infoblog.PostRepository {
	return &postsRepository{db: db}
}
