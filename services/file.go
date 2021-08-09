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

const (
	FileTypeNull = iota
	FileTypePost
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
		UUID: fileUUID,
	}, nil
}

func (f *File) SaveDB(ctx context.Context, file *infoblog.File) error {
	id, err := f.fileRep.Create(ctx, file)
	file.ID = id

	return err
}

func (f *File) GetByIDs(ctx context.Context, ids []string) ([]infoblog.File, error) {
	return f.fileRep.FindByIDs(ctx, ids)
}

func (f *File) UpdatePostUUID(ctx context.Context, ids []string, p infoblog.Post) error {
	return f.fileRep.UpdatePostUUID(ctx, ids, p)
}

func (f *File) GetFilesPostsIDs(ctx context.Context, postsIDs []string) ([]infoblog.File, error) {
	return f.fileRep.FindByPostsIDs(ctx, postsIDs)
}

func (f *File) GetFile(ctx context.Context, fileUUID string) (infoblog.File, error) {
	fileEnt, err := f.fileRep.Find(ctx, infoblog.File{UUID: fileUUID})

	return fileEnt, err
}
