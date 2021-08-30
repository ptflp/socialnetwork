package services

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/auth"
	"gitlab.com/InfoBlogFriends/server/components"
)

type Services struct {
	AuthService infoblog.AuthService
	// TODO change to interface
	User     *User
	Post     *Post
	File     *File
	Comments *Comments
}

func NewServices(cmps components.Componenter, repositories infoblog.Repositories) *Services {
	var services Services
	comments := NewCommentsService(repositories.Comments, &services)
	file := NewFileService(repositories.Files)
	post := NewPostService(repositories, file, cmps.Decoder(), &services)
	user := NewUserService(repositories, post, cmps, file)

	services.AuthService = auth.NewAuthService(repositories, cmps)
	services.Comments = comments
	services.User = user
	services.Post = post
	services.File = file

	return &services
}
