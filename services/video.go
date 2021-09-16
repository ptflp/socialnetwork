package services

import (
	"context"
	"os"
	"path"
	"strings"
	"time"

	"gitlab.com/InfoBlogFriends/server/components"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"

	"gitlab.com/InfoBlogFriends/server/request"
	"gitlab.com/InfoBlogFriends/server/types"
	"gitlab.com/InfoBlogFriends/server/utils"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Video struct {
	file   *File
	ctx    context.Context
	logger *zap.Logger
	ch     chan infoblog.File
}

const (
	ToMp4Path        = "public/videos/new"
	ConvertedMP4Path = "public/videos/converted"
	ToHLSPath        = "public/videos/tohls"
	ConvertedHLSPath = "public/videos/hls"

	NumJobs           = 13
	WorkerPoolDelay   = 5 * time.Second
	QueryDelayMinutes = 1
)

const (
	StatusNotFound = iota + 1
	StatusConvertedMP4
	StatusConvertedHLS
)

func NewVideoService(ctx context.Context, cmps components.Componenter, services *Services) *Video {
	v := &Video{file: services.File, ctx: ctx, logger: cmps.Logger(), ch: make(chan infoblog.File, NumJobs)}
	go v.WorkerPool()
	return v
}

func (v *Video) UploadVideo(ctx context.Context, formFile FormFile) (request.FileData, error) {
	// 1. save file to filesystem
	u, err := extractUser(ctx)
	if err != nil {
		return request.FileData{}, err
	}

	fileUUID := types.NewNullUUID()

	file, err := v.file.SaveFileSystem(formFile, u, fileUUID, "public", "public/videos/new")
	if err != nil {
		return request.FileData{}, err
	}

	// 3. update file info, save to db
	file.Active = 1
	file.Type = types.TypePost
	file.UserUUID = u.UUID
	file.FileType = types.NewNullInt64(FileTypeVideo)

	err = v.file.SaveDB(ctx, &file)
	if err != nil {
		return request.FileData{}, err
	}

	return request.FileData{
		Link: utils.Link(file),
		UUID: file.UUID.String,
	}, nil
}

func (v *Video) WorkerPool() {
	var err error

	ticker := time.NewTicker(WorkerPoolDelay)
	condition := infoblog.Condition{
		Equal: &sq.Eq{"file_type": FileTypeVideo, "status": nil},
		LimitOffset: &infoblog.LimitOffset{
			Offset: 0,
			Limit:  NumJobs,
		},
		Other: &infoblog.Other{
			Condition: "updated_at < date_sub(now(), interval ? minute)",
			Args:      []interface{}{QueryDelayMinutes},
		},
	}
	for w := 0; w < NumJobs; w++ {
		go v.Worker()
	}

	var videos []infoblog.File
	for {
		select {
		case <-v.ctx.Done():
			return
		case <-ticker.C:
			videos, err = v.file.Listx(v.ctx, condition)
			if err != nil {
				v.logger.Error("retrieve videos error", zap.Error(err))
				continue
			}
			for i := range videos {
				v.ch <- videos[i]
			}
		}
	}
}

func (v *Video) Worker() {
	var video infoblog.File
	var err error
	for {
		select {
		case <-v.ctx.Done():
			return
		case video = <-v.ch:
			hlsDir := path.Join(ConvertedHLSPath, video.UUID.String)
			if _, err = os.Stat(hlsDir); os.IsNotExist(err) {
				video.Status = types.NewNullInt64(StatusConvertedHLS)
				err = v.file.Update(context.Background(), video)
				if err != nil {
					v.logger.Error("video update status not found", zap.Error(err))
				}
				continue
			}
			fileName := strings.Join([]string{video.UUID.String, "mp4"}, ".")
			originalFile := path.Join(ToMp4Path, video.Name)
			if _, err = os.Stat(originalFile); os.IsNotExist(err) {
				video.Status = types.NewNullInt64(StatusNotFound)
				err = v.file.Update(context.Background(), video)
				if err != nil {
					v.logger.Error("video update status not found", zap.Error(err))
				}
				continue
			}

			mp4Path := path.Join(ConvertedMP4Path, fileName)
			if _, err = os.Stat(mp4Path); os.IsNotExist(err) {
				continue
			}
			toHLSPath := path.Join(ToHLSPath, fileName)
			err = os.Rename(mp4Path, toHLSPath)
			if err != nil {
				v.logger.Error("move video err", zap.Error(err))
			}

			video.Status = types.NewNullInt64(StatusConvertedMP4)
			err = v.file.Update(context.Background(), video)
			if err != nil {
				err = os.Rename(toHLSPath, mp4Path)
				if err != nil {
					v.logger.Error("move video err", zap.Error(err))
				}
				v.logger.Error("update video status", zap.Error(err))
			}
		}
	}
}
