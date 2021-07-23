package services

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/auth"
	"gitlab.com/InfoBlogFriends/server/components"
)

type Services struct {
	AuthService infoblog.AuthService
	// TODO change to interface
	User *User
	Post *Post
	File *File
}

func NewServices(cmps components.Componenter, repositories infoblog.Repositories) Services {
	file := NewFileService(repositories.Files)
	post := NewPostService(repositories, file, cmps.Decoder())

	return Services{
		AuthService: auth.NewAuthService(repositories, cmps),
		User:        NewUserService(repositories, post, cmps, file),
		Post:        post,
		File:        file,
	}
}
