package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type ChatPrivateUserRepository struct {
	db *sqlx.DB
	crud
}

func (c *ChatPrivateUserRepository) Create(ctx context.Context, chatPrivateUser infoblog.ChatPrivateUser) error {
	return c.create(ctx, &chatPrivateUser)
}

func (c *ChatPrivateUserRepository) Find(ctx context.Context, chatPrivateUser infoblog.ChatPrivateUser) (infoblog.ChatPrivateUser, error) {
	err := c.find(ctx, &chatPrivateUser, &chatPrivateUser)

	return chatPrivateUser, err
}

func (c *ChatPrivateUserRepository) Update(ctx context.Context, chatPrivateUser infoblog.ChatPrivateUser) error {
	return c.update(ctx, &chatPrivateUser)
}

func (c *ChatPrivateUserRepository) Delete(ctx context.Context, chatPrivateUser infoblog.ChatPrivateUser) error {
	chatPrivateUser.Active = types.NullBool{}
	return c.update(ctx, &chatPrivateUser)
}

func (c *ChatPrivateUserRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.ChatPrivateUser, error) {
	var chatPrivateUser []infoblog.ChatPrivateUser
	err := c.list(ctx, &chatPrivateUser, &infoblog.ChatPrivateUser{}, limit, offset)

	return chatPrivateUser, err
}

func (c *ChatPrivateUserRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.ChatPrivateUser, error) {
	var chatPrivateUser []infoblog.ChatPrivateUser
	err := c.crud.listx(ctx, &chatPrivateUser, infoblog.ChatPrivateUser{}, condition)
	if err != nil {
		return nil, err
	}

	return chatPrivateUser, nil
}

func NewChatPrivateUserRepository(db *sqlx.DB) infoblog.ChatPrivateUsersRepository {
	cr := crud{db: db}
	return &ChatPrivateUserRepository{db: db, crud: cr}
}
