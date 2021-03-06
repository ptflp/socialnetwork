package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	UploadDirectory = "./public"
)

const (
	FileTypeNull = iota
	FileTypeImage
	FileTypeVideo
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

func (f *File) SaveFileSystem(formFile FormFile, user infoblog.User, fileUUID types.NullUUID, args ...string) (infoblog.File, error) {
	uploadDir := UploadDirectory
	if len(args) > 0 {
		uploadDir = args[0]
	}
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		err = os.Mkdir(uploadDir, 0755)
		if err != nil {
			return infoblog.File{}, err
		}
	}

	dir := path.Join(uploadDir, user.UUID.String)
	if len(args) > 1 {
		dir = args[1]
	}

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

	fileName := strings.Join([]string{fileUUID.String, s[len(s)-1]}, ".")
	filePath := path.Join(dir, fileName)

	out, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return infoblog.File{}, err
	}
	defer out.Close()

	mtype, err := mimetype.DetectReader(formFile.File)
	if err != nil {
		return infoblog.File{}, err
	}
	_, err = formFile.File.Seek(0, 0)
	if err != nil {
		return infoblog.File{}, err
	}

	_, err = io.Copy(out, formFile.File)

	if err != nil {
		return infoblog.File{}, err
	}

	return infoblog.File{
		Dir:      dir,
		Name:     fileName,
		UUID:     fileUUID,
		MimeType: types.NewNullString(mtype.String()),
	}, nil
}

func (f *File) SaveDB(ctx context.Context, file *infoblog.File) error {
	_, err := f.fileRep.Create(ctx, file)

	return err
}

func (f *File) GetByIDs(ctx context.Context, ids []string) ([]infoblog.File, error) {
	return f.fileRep.FindByIDs(ctx, ids)
}

func (f *File) UpdatePostUUID(ctx context.Context, ids []string, p infoblog.Post) error {
	return f.fileRep.UpdatePostUUID(ctx, ids, p)
}

func (f *File) UpdateFileType(ctx context.Context, ids []string, file infoblog.File) error {
	uuids := make([]types.NullUUID, 0, len(ids))

	for i := range ids {
		uuids = append(uuids, types.NewNullUUID(ids[i]))
	}
	if len(uuids) < 1 {
		return nil
	}

	return f.fileRep.UpdateFileType(ctx, file, uuids...)
}

func (f *File) GetFilesByPostUUIDs(ctx context.Context, postUUIDs []string) ([]infoblog.File, error) {
	return f.fileRep.FindByPostsIDs(ctx, postUUIDs)
}

func (f *File) GetFile(ctx context.Context, fileUUID string) (infoblog.File, error) {
	fileEnt, err := f.fileRep.Find(ctx, infoblog.File{UUID: types.NewNullUUID(fileUUID)})

	return fileEnt, err
}

func (f *File) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.File, error) {
	return f.fileRep.Listx(ctx, condition)
}

func (f *File) Update(ctx context.Context, file infoblog.File) error {
	return f.fileRep.Update(ctx, file)
}
