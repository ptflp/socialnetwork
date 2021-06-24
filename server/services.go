package server

import (
	"gitlab.com/ptflp/infoblog-server"
	"gitlab.com/ptflp/infoblog-server/respond"
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
}
