package services

import (
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
}

func NewServices(cmps components.Componenter, reps infoblog.Repositories) *Services {
	var services Services
	comments := NewCommentsService(reps.Comments, &services)
	file := NewFileService(reps.Files)
	post := NewPostService(reps, file, cmps.Decoder(), &services)
	user := NewUserService(reps, post, cmps, file)
	moderates := NewModeratesService(reps, &services)
	chats := NewChatService(reps, &services)

	services.AuthService = auth.NewAuthService(reps, cmps)
	services.Comments = comments
	services.User = user
	services.Post = post
	services.File = file
	services.Moderates = moderates
	services.Chats = chats

	return &services
}
