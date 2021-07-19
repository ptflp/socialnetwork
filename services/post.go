package services

import (
	"context"
	"fmt"
	"math/rand"
	"path"
	"strings"
	"time"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"github.com/google/uuid"

	"gitlab.com/InfoBlogFriends/server/request"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Post struct {
	file *File
	post infoblog.PostRepository
	*decoder.Decoder
}

func NewPostService(reps infoblog.Repositories, file *File, d *decoder.Decoder) *Post {
	return &Post{post: reps.Posts, file: file, Decoder: d}
}

func (p *Post) SavePost(ctx context.Context, formFile FormFile, req request.PostCreateReq, u *infoblog.User) (request.PostDataResponse, error) {
	// 1. save file to filesystem
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89) + 10

	fUUID, err := uuid.NewUUID()
	if err != nil {
		return request.PostDataResponse{}, err
	}

	fileUUID := strings.Join([]string{fUUID.String(), fmt.Sprintf("-f%d", id)}, "")

	file, err := p.file.SaveFileSystem(formFile, u.ID, fileUUID)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	// 2. save post to db

	pUUID, err := uuid.NewUUID()
	if err != nil {
		return request.PostDataResponse{}, err
	}
	postUUID := strings.Join([]string{pUUID.String(), fmt.Sprintf("-p%d", id)}, "")

	post := infoblog.Post{
		PostEntity: infoblog.PostEntity{
			Body:     req.Body,
			UserID:   u.ID,
			Type:     1,
			UUID:     postUUID,
			UserUUID: u.UUID,
		},
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
	file.ForeignUUID = postUUID

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
		UUID: post.UUID,
		Body: post.Body,
		Files: []request.PostFileData{
			{
				Link: "/" + path.Join(file.Dir, file.Name),
				UUID: file.UUID,
			},
		},
		User: request.UserData{
			UUID:       u.UUID,
			Name:       "",
			SecondName: "",
		},
		Counts: request.PostCountData{
			Likes:    0,
			Comments: 0,
		},
	}, nil
}

func (p *Post) FeedRecent(ctx context.Context, req request.PostsFeedReq) (request.PostsFeedData, error) {
	posts, postIDIndexMap, postsIDs, err := p.post.FindAllRecent(ctx, req.Limit, req.Offset)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	if len(posts) < 1 {
		return request.PostsFeedData{}, nil
	}
	files, err := p.file.GetFilesPostsIDs(ctx, postsIDs)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	count, err := p.post.CountRecent(ctx)
	if err != nil {
		return request.PostsFeedData{}, err
	}

	for i := range files {
		id := postIDIndexMap[files[i].ForeignID]
		posts[id].Files = append(posts[id].Files, files[i])
	}

	postDataRes := make([]request.PostDataResponse, 0, req.Limit)

	for i := range posts {
		postsFileData := make([]request.PostFileData, 0, req.Limit)
		for j := range posts[i].Files {
			postsFileData = append(postsFileData, request.PostFileData{
				Link: "/" + path.Join(posts[i].Files[j].Dir, posts[i].Files[j].Name),
				UUID: posts[i].Files[j].UUID,
			})
		}

		userData := request.UserData{}

		err = p.MapStructs(&userData, &posts[i].User)
		if err != nil {
			return request.PostsFeedData{}, err
		}

		pdr := request.PostDataResponse{
			Files: postsFileData,
			User:  userData,
		}

		err = p.MapStructs(&pdr, &posts[i].PostEntity)
		if err != nil {
			return request.PostsFeedData{}, err
		}

		postDataRes = append(postDataRes, pdr)
	}

	return request.PostsFeedData{
		Count: count,
		Posts: postDataRes,
	}, nil
}

func (p *Post) FeedMy(ctx context.Context, u infoblog.User, req request.PostsFeedReq) (request.PostsFeedData, error) {
	posts, postIDIndexMap, postsIDs, err := p.post.FindAll(ctx, u, req.Limit, req.Offset)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	if len(posts) < 1 {
		return request.PostsFeedData{}, nil
	}
	files, err := p.file.GetFilesPostsIDs(ctx, postsIDs)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	count, err := p.post.CountRecent(ctx)
	if err != nil {
		return request.PostsFeedData{}, err
	}

	for i := range files {
		id := postIDIndexMap[files[i].ForeignID]
		posts[id].Files = append(posts[id].Files, files[i])
	}

	postDataRes := make([]request.PostDataResponse, 0, req.Limit)

	for i := range posts {
		postsFileData := make([]request.PostFileData, 0, req.Limit)
		for j := range posts[i].Files {
			postsFileData = append(postsFileData, request.PostFileData{
				Link: "/" + path.Join(posts[i].Files[j].Dir, posts[i].Files[j].Name),
				UUID: posts[i].Files[j].UUID,
			})
		}

		postDataRes = append(postDataRes, request.PostDataResponse{
			UUID:  posts[i].UUID,
			Body:  posts[i].Body,
			Files: postsFileData,
			User: request.UserData{
				UUID:       posts[i].User.UUID,
				Name:       posts[i].User.Name.String,
				SecondName: posts[i].User.SecondName.String,
			},
			Counts: request.PostCountData{
				Likes:    0,
				Comments: 0,
			},
		})
	}

	return request.PostsFeedData{
		Count: count,
		Posts: postDataRes,
	}, nil
}

func (p *Post) savePostDB(ctx context.Context, pst *infoblog.Post) error {
	_, err := p.post.Create(ctx, *pst)
	if err != nil {
		return err
	}

	post, err := p.post.Find(ctx, *pst)
	if err != nil {
		return err
	}
	*pst = post

	return err
}
