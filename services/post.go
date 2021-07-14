package services

import (
	"context"
	"path"

	"gitlab.com/InfoBlogFriends/server/request"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Post struct {
	file *File
	post infoblog.PostRepository
}

func NewPostService(reps infoblog.Repositories, file *File) *Post {
	return &Post{post: reps.Posts, file: file}
}

func (p *Post) SavePost(ctx context.Context, formFile FormFile, req request.PostCreateReq, u *infoblog.User) (request.PostDataResponse, error) {
	// 1. save file to filesystem
	file, err := p.file.SaveFileSystem(formFile, u.ID)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	// 2. save post to db
	post := infoblog.Post{
		Body:   req.Body,
		UserID: u.ID,
		Type:   1,
	}
	err = p.savePostDB(ctx, &post)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	// 3. update file info, save to db
	file.ForeignID = post.ID
	file.Active = 1
	file.Type = 1
	file.UserID = u.ID

	err = p.file.SaveDB(ctx, &file)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	// 4. update post and activate
	post.FileID = file.ID
	post.Active = 1
	err = p.post.Update(ctx, post)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	return request.PostDataResponse{
		ID:    post.ID,
		Body:  post.Body,
		Files: []string{"/" + path.Join(file.Dir, file.Name)},
		User: request.UserData{
			ID:         u.ID,
			Name:       "",
			SecondName: "",
		},
		Counts: request.PostCountData{
			Likes:    0,
			Comments: 0,
		},
	}, nil
}

func (p *Post) savePostDB(ctx context.Context, pst *infoblog.Post) error {
	id, err := p.post.Create(ctx, *pst)
	if err != nil {
		return err
	}

	post, err := p.post.Find(ctx, id)
	if err != nil {
		return err
	}
	*pst = post

	return err
}
