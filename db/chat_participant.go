package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type ChatParticipantRepository struct {
	db *sqlx.DB
	crud
}

func (c *ChatParticipantRepository) Create(ctx context.Context, chatMessage infoblog.ChatParticipant) error {
	return c.create(ctx, &chatMessage)
}

func (c *ChatParticipantRepository) Find(ctx context.Context, chatMessage infoblog.ChatParticipant) (infoblog.ChatParticipant, error) {
	err := c.find(ctx, &chatMessage, &chatMessage)

	return chatMessage, err
}

func (c *ChatParticipantRepository) Update(ctx context.Context, chatMessage infoblog.ChatParticipant) error {
	return c.update(ctx, &chatMessage)
}

func (c *ChatParticipantRepository) Delete(ctx context.Context, chatMessage infoblog.ChatParticipant) error {
	chatMessage.Active = types.NullBool{}
	return c.update(ctx, &chatMessage)
}

func (c *ChatParticipantRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.ChatParticipant, error) {
	var chatParticipant []infoblog.ChatParticipant
	err := c.list(ctx, &chatParticipant, &infoblog.ChatParticipant{}, limit, offset)

	return chatParticipant, err
}

func NewChatParticipantRepository(db *sqlx.DB) infoblog.ChatParticipantRepository {
	cr := crud{db: db}
	return &ChatParticipantRepository{db: db, crud: cr}
}
