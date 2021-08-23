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

func (c *CommentsRepository) Create(ctx context.Context, comments infoblog.Comments) error {
	return c.create(ctx, &comments)
}

func (c *CommentsRepository) Find(ctx context.Context, comments infoblog.Comments) (infoblog.Comments, error) {
	err := c.find(ctx, &comments, &comments)

	return comments, err
}

func (c *CommentsRepository) Update(ctx context.Context, comments infoblog.Comments) error {
	return c.update(ctx, &comments)
}

func (c *CommentsRepository) Delete(ctx context.Context, comments infoblog.Comments) error {
	comments.Active = types.NullBool{}
	return c.update(ctx, &comments)
}

func (c *CommentsRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Comments, error) {
	var chat []infoblog.Comments
	err := c.list(ctx, &chat, &infoblog.Comments{}, limit, offset)

	return chat, err
}

func NewCommentsRepository(db *sqlx.DB) infoblog.CommentsRepository {
	cr := crud{db: db}
	return &CommentsRepository{db: db, crud: cr}
}
