package server

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/email"
	"gitlab.com/InfoBlogFriends/server/respond"
	"gitlab.com/InfoBlogFriends/server/service"
	"gitlab.com/InfoBlogFriends/server/session"
	"go.uber.org/zap"
)

type Services struct {
	AuthService infoblog.AuthService
	User        *service.User
	Post        *service.Post
	File        *service.File
}

type Components struct {
	UserRepository infoblog.UserRepository
	Logger         *zap.Logger
	Responder      respond.Responder
	LogLevel       zap.AtomicLevel
	JWTKeys        *session.JWTKeys
	Email          *email.Client
}
