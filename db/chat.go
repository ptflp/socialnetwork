package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type ChatRepository struct {
	db *sqlx.DB
	crud
}

func (c *ChatRepository) Create(ctx context.Context, chatMessage infoblog.Chat) error {
	return c.create(ctx, &chatMessage)
}

func (c *ChatRepository) Find(ctx context.Context, chatMessage infoblog.Chat) (infoblog.Chat, error) {
	err := c.find(ctx, &chatMessage, &chatMessage)

	return chatMessage, err
}

func (c *ChatRepository) Update(ctx context.Context, chatMessage infoblog.Chat) error {
	return c.update(ctx, &chatMessage)
}

func (c *ChatRepository) Delete(ctx context.Context, chatMessage infoblog.Chat) error {
	chatMessage.Active = types.NullBool{}
	return c.update(ctx, &chatMessage)
}

func (c *ChatRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Chat, error) {
	var chat []infoblog.Chat
	err := c.list(ctx, &chat, &infoblog.Chat{}, limit, offset)

	return chat, err
}

func NewChatRepository(db *sqlx.DB) infoblog.ChatRepository {
	cr := crud{db: db}
	return &ChatRepository{db: db, crud: cr}
}
