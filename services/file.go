package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	UploadDirectory = "./public"
)

type File struct {
	fileRep infoblog.FileRepository
}

type FormFile struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

func NewFileService(fileRep infoblog.FileRepository) *File {
	return &File{fileRep: fileRep}
}

func (f *File) SaveFileSystem(formFile FormFile, uid int64, fileUUID string) (infoblog.File, error) {
	if _, err := os.Stat(UploadDirectory); os.IsNotExist(err) {
		err = os.Mkdir(UploadDirectory, 0755)
		if err != nil {
			return infoblog.File{}, err
		}
	}

	dir := path.Join(UploadDirectory, strconv.Itoa(int(uid)))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return infoblog.File{}, err
		}
	}

	if !strings.Contains(formFile.FileHeader.Filename, ".") {
		return infoblog.File{}, fmt.Errorf("filename without extension %s", formFile.FileHeader.Filename)
	}

	s := strings.Split(formFile.FileHeader.Filename, ".")

	fileName := strings.Join([]string{fileUUID, s[len(s)-1]}, ".")
	filePath := path.Join(dir, fileName)

	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return infoblog.File{}, err
	}
	defer out.Close()

	_, err = io.Copy(out, formFile.File)

	if err != nil {
		return infoblog.File{}, err
	}

	return infoblog.File{
		Dir:  dir,
		Name: fileName,
	}, nil
}

func (f *File) SaveDB(ctx context.Context, file *infoblog.File) error {
	id, err := f.fileRep.Create(ctx, file)
	file.ID = id

	return err
}

func (f *File) InitPostFile(ctx context.Context, file *infoblog.File) error {
	fileEntity, err := f.fileRep.Find(ctx, file.ID)
	if err != nil {
		return err
	}

	fileEntity.Active = 1
	fileEntity.Type = 1

	return f.fileRep.Update(ctx, *file)
}
