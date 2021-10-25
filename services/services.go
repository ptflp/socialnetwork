package services

import (
	"context"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/auth"
	"gitlab.com/InfoBlogFriends/server/components"
)

type Services struct {
	AuthService infoblog.AuthService
	// TODO change to interface
	User      *User
	Post      *Post
	File      *File
	Comments  *Comments
	Moderates *Moderates
	Chats     *Chats
	Video     *Video
	Event     *Event
}

func NewServices(ctx context.Context, cmps components.Componenter, reps infoblog.Repositories) *Services {
	var services Services
	event := NewEventService(ctx, cmps, reps)
	services.Event = event

	comments := NewCommentsService(reps.Comments, &services)
	services.Comments = comments
	file := NewFileService(reps.Files)
	services.File = file
	post := NewPostService(reps, file, cmps.Decoder(), &services)
	services.Post = post
	user := NewUserService(reps, post, cmps, file, &services)
	services.User = user
	moderates := NewModeratesService(reps, &services)
	services.Moderates = moderates
	chats := NewChatService(reps, &services)
	services.Chats = chats
	video := NewVideoService(ctx, cmps, &services)
	services.Video = video

	services.AuthService = auth.NewAuthService(reps, cmps)

	return &services
}
