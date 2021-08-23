package services

import (
	"context"
	"errors"

	"gitlab.com/InfoBlogFriends/server/types"

	"gitlab.com/InfoBlogFriends/server/decoder"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/request"
)

type Chat struct {
	*decoder.Decoder
	chatRepository infoblog.ChatRepository
}

func NewChatService(rs infoblog.Repositories) *Chat {
	return &Chat{chatRepository: rs.Chats, Decoder: decoder.NewDecoder()}
}

func (c *Chat) saveChatDB(ctx context.Context, cht *infoblog.Chat) error {
	err := c.chatRepository.Create(ctx, *cht)
	if err != nil {
		return err
	}

	chat, err := c.chatRepository.Find(ctx, *cht)
	if err != nil {
		return err
	}
	*cht = chat

	return err
}

func (c *Chat) GetChat(ctx context.Context) (request.ChatData, error) {
	chat, err := extractChat(ctx)
	if err != nil {
		return request.ChatData{}, err
	}

	chat, err = c.chatRepository.Find(ctx, chat)
	if err != nil {
		return request.ChatData{}, err
	}
	chatData := request.ChatData{}
	err = c.MapStructs(&chatData, &chat)
	if err != nil {
		return request.ChatData{}, err
	}

	if err != nil {
		return request.ChatData{}, err
	}

	return chatData, nil
}

func (c *Chat) UpdateChat(ctx context.Context, profileUpdateReq request.ChatUpdateReq, chat infoblog.Chat) (request.ChatData, error) {
	chat, err := c.chatRepository.Find(ctx, chat)
	if err != nil {
		return request.ChatData{}, err
	}

	err = c.MapStructs(&chat, &profileUpdateReq)
	if err != nil {
		return request.ChatData{}, err
	}

	chatData := request.ChatData{}
	err = c.MapStructs(&chatData, &chat)
	if err != nil {
		return request.ChatData{}, err
	}

	return chatData, nil
}

func (c *Chat) Get(ctx context.Context, req request.ChatIDRequest) (request.ChatData, error) {
	chat := infoblog.Chat{}
	var err error
	if req.UUID != nil {
		chat.UUID = types.NewNullUUID(*req.UUID)
		chat, err = c.chatRepository.Find(ctx, chat)
		if err != nil {
			return request.ChatData{}, err
		}
	}

	chatData := request.ChatData{}
	err = c.MapStructs(&chatData, &chat)
	if err != nil {
		return request.ChatData{}, err
	}

	return chatData, nil
}

func extractChat(ctx context.Context) (infoblog.Chat, error) {
	c, ok := ctx.Value(types.Chat{}).(*infoblog.Chat)
	if !ok {
		return infoblog.Chat{}, errors.New("type assertion to chat err")
	}

	return *c, nil
}
