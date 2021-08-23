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

func (c *ChatMessagesRepository) Create(ctx context.Context, chatMessage infoblog.ChatMessages) error {
	return c.create(ctx, &chatMessage)
}

func (c *ChatMessagesRepository) Find(ctx context.Context, chatMessage infoblog.ChatMessages) (infoblog.ChatMessages, error) {
	err := c.find(ctx, &chatMessage, &chatMessage)

	return chatMessage, err
}

func (c *ChatMessagesRepository) Update(ctx context.Context, chatMessage infoblog.ChatMessages) error {
	return c.update(ctx, &chatMessage)
}

func (c *ChatMessagesRepository) Delete(ctx context.Context, chatMessage infoblog.ChatMessages) error {
	chatMessage.Active = types.NullBool{}
	return c.update(ctx, &chatMessage)
}

func (c *ChatMessagesRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.ChatMessages, error) {
	var chatMessages []infoblog.ChatMessages
	err := c.list(ctx, &chatMessages, &infoblog.ChatMessages{}, limit, offset)

	return chatMessages, err
}

func NewChatMessagesRepository(db *sqlx.DB) infoblog.ChatMessagesRepository {
	cr := crud{db: db}
	return &ChatMessagesRepository{db: db, crud: cr}
}
