package utils

import (
	"path"

	"gitlab.com/InfoBlogFriends/server/request"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

func Link(file infoblog.File) string {
	return path.Join("/file", file.UUID.String)
}

func PrivateLink(file request.PostFileData) string {
	return path.Join("/file", file.UUID)
}

func PublicLink(file request.PostFileData) string {
	return path.Join("/", file.Dir, file.Name)
}
