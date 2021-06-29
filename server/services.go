package server

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/respond"
	"gitlab.com/InfoBlogFriends/server/session"
	"go.uber.org/zap"
)

type Services struct {
	AuthService infoblog.AuthService
}

type HandlerComponents struct {
	UserRepository infoblog.UserRepository
	Logger         *zap.Logger
	Responder      respond.Responder
	LogLevel       zap.AtomicLevel
	JWTKeys        *session.JWTKeys
}
