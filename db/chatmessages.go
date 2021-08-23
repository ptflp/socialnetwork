package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type ChatMessagesRepository struct {
	db *sqlx.DB
	crud
}

func (c *ChatMessagesRepository) Create(ctx context.Context, chatMessages infoblog.ChatMessages) error {
	return c.create(ctx, &chatMessages)
}

func (c *ChatMessagesRepository) Find(ctx context.Context, chatMessages infoblog.ChatMessages) (infoblog.ChatMessages, error) {
	err := c.find(ctx, &chatMessages, &chatMessages)

	return chatMessages, err
}

func (c *ChatMessagesRepository) Update(ctx context.Context, chatMessages infoblog.ChatMessages) error {
	return c.update(ctx, &chatMessages)
}

func (c *ChatMessagesRepository) Delete(ctx context.Context, chatMessages infoblog.ChatMessages) error {
	chatMessages.Active = types.NullBool{}
	return c.update(ctx, &chatMessages)
}

func (c *ChatMessagesRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.ChatMessages, error) {
	var chatMessagess []infoblog.ChatMessages
	err := c.list(ctx, &chatMessagess, &infoblog.ChatMessages{}, limit, offset)

	return chatMessagess, err
}

func NewChatMessagesRepository(db *sqlx.DB) infoblog.ChatMessagesRepository {
	cr := crud{db: db}
	return &ChatMessagesRepository{db: db, crud: cr}
}
