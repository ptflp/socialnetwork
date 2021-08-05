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
	services *Services
	file     *File
	post     infoblog.PostRepository
	like     infoblog.LikeRepository
	*decoder.Decoder
}

func NewPostService(reps infoblog.Repositories, file *File, d *decoder.Decoder, services *Services) *Post {
	return &Post{post: reps.Posts, file: file, Decoder: d, like: reps.Likes, services: services}
}

func (p *Post) SaveFile(ctx context.Context, formFile FormFile) (request.PostFileData, error) {
	// 1. save file to filesystem
	u, err := extractUser(ctx)
	if err != nil {
		return request.PostFileData{}, err
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89) + 10

	fUUID, err := uuid.NewUUID()
	if err != nil {
		return request.PostFileData{}, err
	}

	fileUUID := strings.Join([]string{fUUID.String(), fmt.Sprintf("-f%d", id)}, "")

	file, err := p.file.SaveFileSystem(formFile, u.ID, fileUUID)
	if err != nil {
		return request.PostFileData{}, err
	}

	// 2. save post to db

	// 3. update file info, save to db
	file.Active = 1
	file.Type = 1
	file.UserUUID = u.UUID

	err = p.file.SaveDB(ctx, &file)
	if err != nil {
		return request.PostFileData{}, err
	}

	return request.PostFileData{
		Link: "/" + path.Join(file.Dir, file.Name),
		UUID: file.UUID,
	}, nil
}

func (p *Post) SavePost(ctx context.Context, req request.PostCreateReq) (request.PostDataResponse, error) {
	u, err := extractUser(ctx)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	var price infoblog.NullFloat64
	if req.PostType == 2 {
		if req.Price == nil {
			return request.PostDataResponse{}, fmt.Errorf("bad price %s", req.Price)
		}
		price = infoblog.NewNullFloat64(*req.Price)
	}

	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89) + 10
	pUUID, err := uuid.NewUUID()
	if err != nil {
		return request.PostDataResponse{}, err
	}
	postUUID := strings.Join([]string{pUUID.String(), fmt.Sprintf("-p%d", id)}, "")
	post := infoblog.Post{
		PostEntity: infoblog.PostEntity{
			Body:     req.Description,
			UserID:   u.ID,
			Type:     req.PostType,
			UUID:     postUUID,
			UserUUID: u.UUID,
			Active:   1,
			Price:    price,
		},
	}
	err = p.savePostDB(ctx, &post)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	err = p.file.UpdatePostUUID(ctx, req.FilesID, post)
	// 1. save file to filesystem
	filesRaw, err := p.file.GetByIDs(ctx, req.FilesID)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	files := make([]request.PostFileData, 0, len(filesRaw))

	for i := range filesRaw {
		file := request.PostFileData{
			Link: "/" + path.Join(filesRaw[i].Dir, filesRaw[i].Name),
			UUID: filesRaw[i].UUID,
		}
		files = append(files, file)
	}

	postDataRes := request.PostDataResponse{}
	postDataRes.Files = files

	err = p.MapStructs(&postDataRes, &post.PostEntity)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	err = p.MapStructs(&postDataRes.User, &u)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	// 4. update post and activate
	return postDataRes, nil
}

func (p *Post) Update(ctx context.Context, req request.PostUpdateReq) error {
	u, err := extractUser(ctx)
	var post infoblog.Post
	post.UUID = req.UUID
	post, err = p.post.Find(ctx, post)
	if err != nil {
		return err
	}

	if post.UserUUID != u.UUID {
		return fmt.Errorf("permission denied")
	}

	if req.Price != nil {
		post.Price = infoblog.NewNullFloat64(*req.Price)
	}
	post.Body = req.Body

	return p.post.Update(ctx, post)
}

func (p *Post) Delete(ctx context.Context, req request.PostUUIDReq) error {
	u, err := extractUser(ctx)
	var post infoblog.Post
	post.UUID = req.UUID
	post, err = p.post.Find(ctx, post)
	if err != nil {
		return err
	}

	if post.UserUUID != u.UUID {
		return fmt.Errorf("permission denied")
	}

	return p.post.Delete(ctx, post)
}

func (p *Post) Get(ctx context.Context, req request.PostUUIDReq) (request.PostDataResponse, error) {
	var err error
	postDataRes := request.PostDataResponse{}
	post := infoblog.Post{}
	post.UUID = req.UUID
	post, err = p.post.Find(ctx, post)
	if err != nil {
		return postDataRes, err
	}

	filesRaw, err := p.file.GetFilesPostsIDs(ctx, []string{req.UUID})
	if err != nil {
		return postDataRes, err
	}

	files := make([]request.PostFileData, 0, len(filesRaw))

	for i := range filesRaw {
		file := request.PostFileData{
			Link: "/" + path.Join(filesRaw[i].Dir, filesRaw[i].Name),
			UUID: filesRaw[i].UUID,
		}
		files = append(files, file)
	}
	err = p.MapStructs(&postDataRes, &post.PostEntity)
	postDataRes.Files = files
	var u infoblog.User
	u.UUID = post.UserUUID
	u, err = p.services.User.userRepository.Find(ctx, u)
	err = p.MapStructs(&postDataRes.User, &u)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	return postDataRes, nil
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
		id := postIDIndexMap[files[i].ForeignUUID]
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

		userData, err = p.services.User.GetUserData(posts[i].User)
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

		pdr.Counts.Likes, err = p.like.CountByPost(ctx, infoblog.Like{Type: 1, ForeignUUID: pdr.UUID})
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

func (p *Post) FeedByUser(ctx context.Context, req request.PostsFeedUserReq) (request.PostsFeedData, error) {
	u := infoblog.User{UUID: req.UUID}
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
	count, err := p.post.CountByUser(ctx, u)
	if err != nil {
		return request.PostsFeedData{}, err
	}

	for i := range files {
		id := postIDIndexMap[files[i].ForeignUUID]
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

		userData, err = p.services.User.GetUserData(posts[i].User)
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

		pdr.Counts.Likes, err = p.like.CountByPost(ctx, infoblog.Like{Type: 1, ForeignUUID: pdr.UUID})
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

func (p *Post) CountByUser(ctx context.Context, user infoblog.User) (int64, error) {
	return p.post.CountByUser(ctx, user)
}

func (p *Post) Like(ctx context.Context, req request.PostUUIDReq) error {
	u, ok := ctx.Value("user").(*infoblog.User)
	if !ok {
		return fmt.Errorf("get user from request context err")
	}

	post, err := p.post.Find(ctx, infoblog.Post{
		PostEntity: infoblog.PostEntity{UUID: req.UUID},
	})

	if err != nil {
		return err
	}
	like := infoblog.Like{
		Type:        1,
		ForeignUUID: post.UUID,
		UserUUID:    post.UserUUID,
		LikerUUID:   u.UUID,
	}

	likeFound, err := p.like.Find(ctx, &like)
	if err != nil {
		like.Active = infoblog.NewNullBool(true)
		return p.like.Upsert(ctx, like)
	}
	likeFound.Active = infoblog.NewNullBool(!likeFound.Active.Bool)

	return p.like.Upsert(ctx, likeFound)
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
