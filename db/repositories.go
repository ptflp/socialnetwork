package db

import (
	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/components"
	"gitlab.com/InfoBlogFriends/server/migration"
	"go.uber.org/zap"
)

func NewRepositories(cmps components.Componenter) infoblog.Repositories {

	mainDB, err := NewDB(cmps.Logger(), cmps.Config().DB)
	if err != nil {
		cmps.Logger().Fatal("db initialization error", zap.Error(err))
	}
	migrator := migration.NewMigrator(mainDB)
	err = migrator.Migrate()
	if err != nil {
		cmps.Logger().Fatal("error on migration apply", zap.Error(err))
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
