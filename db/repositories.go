package db

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/components"
	"go.uber.org/zap"
)

func NewRepositories(cmps components.Componenter) infoblog.Repositories {

	mainDB, err := NewDB(cmps.Logger(), cmps.Config().DB)
	if err != nil {
		cmps.Logger().Fatal("db initialization error", zap.Error(err))
	}

	r := infoblog.Repositories{
		Files:       NewFilesRepository(mainDB),
		Posts:       NewPostsRepository(mainDB),
		Users:       NewUserRepository(mainDB),
		Subscribers: NewSubscribeRepository(mainDB),
		Likes:       NewLikesRepository(mainDB),
	}

	return r
}
