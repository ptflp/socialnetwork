package services

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"gitlab.com/InfoBlogFriends/server/decoder"
	"gitlab.com/InfoBlogFriends/server/request"
	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Comments struct {
	commentsRep infoblog.CommentsRepository
	Services    *Services
	*decoder.Decoder
}

func NewCommentsService(commentsRep infoblog.CommentsRepository, services *Services) *Comments {
	return &Comments{commentsRep: commentsRep, Decoder: decoder.NewDecoder(), Services: services}
}

func (c *Comments) CreatePostComment(ctx context.Context, commentReq request.CommentCreateReq) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}
	comment := infoblog.Comment{
		UUID:        types.NewNullUUID(),
		Type:        types.TypePost,
		UserUUID:    user.UUID,
		Active:      types.NewNullBool(true),
		ForeignUUID: types.NewNullUUID(commentReq.ForeignUUID),
		Body:        types.NewNullString(commentReq.Body),
	}

	err = c.commentsRep.Create(ctx, comment)
	if err != nil {
		return err
	}
	go c.Services.Post.UpdateCounters(ctx, commentReq.ForeignUUID)

	return nil
}

func (c *Comments) GetPostComments(ctx context.Context, commentReq request.PostUUIDReq) ([]request.CommentData, error) {
	condition := infoblog.Condition{
		Equal: &sq.Eq{"type": types.TypePost, "foreign_uuid": types.NewNullUUID(commentReq.UUID), "active": true},
		Order: &infoblog.Order{
			Field: "created_at",
		},
	}
	comments, err := c.commentsRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}
	if len(comments) < 1 {
		return []request.CommentData{}, nil
	}
	commentsData := make([]request.CommentData, len(comments))
	err = c.MapStructs(&commentsData, &comments)
	if err != nil {
		return nil, err
	}

	usersUUIDs := make([]interface{}, 0, len(commentsData))
	mapComment := make(map[string][]*request.CommentData, len(commentsData))
	for i := range commentsData {
		usersUUIDs = append(usersUUIDs, commentsData[i].UserUUID)
		mapComment[commentsData[i].UserUUID.String] = append(mapComment[commentsData[i].UserUUID.String], &commentsData[i])
	}
	condition = infoblog.Condition{
		Equal: &sq.Eq{"active": true},
		In: &infoblog.In{
			Field: "uuid",
			Args:  usersUUIDs,
		},
	}

	users, err := c.Services.User.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	usersData := make([]request.UserData, len(users))

	err = c.MapStructs(&usersData, &users)
	if err != nil {
		return nil, err
	}

	for i := range usersData {
		for j := range mapComment[usersData[i].UUID.String] {
			mapComment[usersData[i].UUID.String][j].User = usersData[i]
		}
	}

	return commentsData, nil
}
