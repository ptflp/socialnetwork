package db

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type CommentsRepository struct {
	db *sqlx.DB
	crud
}

func (c *CommentsRepository) Create(ctx context.Context, comments infoblog.Comment) error {
	return c.create(ctx, &comments)
}

func (c *CommentsRepository) Find(ctx context.Context, comments infoblog.Comment) (infoblog.Comment, error) {
	err := c.find(ctx, &comments, &comments)

	return comments, err
}

func (c *CommentsRepository) Update(ctx context.Context, comments infoblog.Comment) error {
	return c.update(ctx, &comments)
}

func (c *CommentsRepository) Delete(ctx context.Context, comments infoblog.Comment) error {
	comments.Active = types.NullBool{}
	return c.update(ctx, &comments)
}

func (c *CommentsRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Comment, error) {
	var chat []infoblog.Comment
	err := c.list(ctx, &chat, &infoblog.Comment{}, limit, offset)

	return chat, err
}

func (c *CommentsRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.Comment, error) {
	var comments []infoblog.Comment
	err := c.crud.listx(ctx, &comments, infoblog.Comment{}, condition)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentsRepository) GetCount(ctx context.Context, condition infoblog.Condition) (uint64, error) {
	return c.crud.getCount(ctx, infoblog.Comment{}, condition)
}

func NewCommentsRepository(db *sqlx.DB) infoblog.CommentsRepository {
	cr := crud{db: db}
	return &CommentsRepository{db: db, crud: cr}
}
