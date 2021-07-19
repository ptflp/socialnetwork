package db

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/components"
	"go.uber.org/zap"
)

func NewRepositories(cmps components.Componenter) infoblog.Repositories {

	database, err := NewDB(cmps.Logger(), cmps.Config().DB)
	if err != nil {
		cmps.Logger().Fatal("db initialization error", zap.Error(err))
	}

	r := infoblog.Repositories{
		Files:       NewFilesRepository(database),
		Posts:       NewPostsRepository(database),
		Users:       NewUserRepository(database),
		Subscribers: NewSubscribeRepository(database),
	}

	return r
}
