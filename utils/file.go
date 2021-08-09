package utils

import (
	"path"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

func Link(file infoblog.File) string {
	return path.Join("/file", file.UUID)
}
