package services

import (
	"context"
	"fmt"

	"gitlab.com/InfoBlogFriends/server/types"
	"gitlab.com/InfoBlogFriends/server/utils"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

const (
	PostTypeFree = iota
	PostTypeSubscribe
	PostTypeForPrice
)

type Post struct {
	services   *Services
	file       *File
	subscribes infoblog.SubscriberRepository
	post       infoblog.PostRepository
	like       infoblog.LikeRepository
	*decoder.Decoder
}

func NewPostService(reps infoblog.Repositories, file *File, d *decoder.Decoder, services *Services) *Post {
	return &Post{post: reps.Posts, file: file, Decoder: d, like: reps.Likes, services: services, subscribes: reps.Subscribers}
}

func (p *Post) SaveFile(ctx context.Context, formFile FormFile) (request.PostFileData, error) {
	// 1. save file to filesystem
	u, err := extractUser(ctx)
	if err != nil {
		return request.PostFileData{}, err
	}

	fileUUID := types.NewNullUUID()

	file, err := p.file.SaveFileSystem(formFile, u, fileUUID)
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
		Link: utils.Link(file),
		UUID: file.UUID.String,
	}, nil
}

func (p *Post) SavePost(ctx context.Context, req request.PostCreateReq) (request.PostDataResponse, error) {
	u, err := extractUser(ctx)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	var price types.NullFloat64
	if req.PostType == PostTypeForPrice {
		if req.Price == nil {
			return request.PostDataResponse{}, fmt.Errorf("bad price %d", req.Price)
		}
		price = types.NewNullFloat64(*req.Price)
	}

	post := infoblog.Post{
		PostEntity: infoblog.PostEntity{
			Body:     req.Description,
			Type:     req.PostType,
			UUID:     types.NewNullUUID(),
			UserUUID: u.UUID,
			Active:   types.NewNullBool(true),
			Price:    price,
		},
	}
	err = p.savePostDB(ctx, &post)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	if len(req.FilesID) < 1 {
		return request.PostDataResponse{}, fmt.Errorf("no files present")
	}

	_, err = p.file.fileRep.Find(ctx, infoblog.File{UUID: types.NewNullUUID(req.FilesID[0])})
	if err != nil {
		return request.PostDataResponse{}, err
	}

	err = p.file.UpdatePostUUID(ctx, req.FilesID, post)
	if err != nil {
		return request.PostDataResponse{}, err
	}
	// 1. save file to filesystem
	filesRaw, err := p.file.GetByIDs(ctx, req.FilesID)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	files := make([]request.PostFileData, 0, len(filesRaw))

	for i := range filesRaw {
		file := request.PostFileData{
			Link: utils.Link(filesRaw[i]),
			UUID: filesRaw[i].UUID.String,
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
	if err != nil {
		return err
	}
	var post infoblog.Post
	post.UUID = types.NewNullUUID(req.UUID)
	post, err = p.post.Find(ctx, post)
	if err != nil {
		return err
	}

	if post.UserUUID.String != u.UUID.String {
		return fmt.Errorf("permission denied")
	}

	if req.Price != nil {
		post.Price = types.NewNullFloat64(*req.Price)
	}
	post.Body = req.Body

	return p.post.Update(ctx, post)
}

func (p *Post) Delete(ctx context.Context, req request.PostUUIDReq) error {
	u, err := extractUser(ctx)
	if err != nil {
		return err
	}
	var post infoblog.Post
	post.UUID = types.NewNullUUID(req.UUID)
	post, err = p.post.Find(ctx, post)
	if err != nil {
		return err
	}

	if post.UserUUID.String != u.UUID.String {
		return fmt.Errorf("permission denied")
	}

	return p.post.Delete(ctx, post)
}

func (p *Post) Get(ctx context.Context, req request.PostUUIDReq) (request.PostDataResponse, error) {
	var err error
	postDataRes := request.PostDataResponse{}
	post := infoblog.Post{}
	post.UUID = types.NewNullUUID(req.UUID)
	post, err = p.post.Find(ctx, post)
	if err != nil {
		return postDataRes, err
	}

	filesRaw, err := p.file.GetFilesByPostUUIDs(ctx, []string{req.UUID})
	if err != nil {
		return postDataRes, err
	}

	files := make([]request.PostFileData, 0, len(filesRaw))

	for i := range filesRaw {
		file := request.PostFileData{
			Link: utils.Link(filesRaw[i]),
			UUID: filesRaw[i].UUID.String,
		}
		files = append(files, file)
	}
	err = p.MapStructs(&postDataRes, &post.PostEntity)
	if err != nil {
		return request.PostDataResponse{}, err
	}
	postDataRes.Files = files
	var u infoblog.User
	u.UUID = post.UserUUID
	u, err = p.services.User.userRepository.Find(ctx, u)
	if err != nil {
		return request.PostDataResponse{}, err
	}
	err = p.MapStructs(&postDataRes.User, &u)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	return postDataRes, nil
}

func (p *Post) FeedRecent(ctx context.Context, req request.LimitOffsetReq) (request.PostsFeedData, error) {
	posts, postIDIndexMap, postsIDs, err := p.post.FindAllRecent(ctx, req.Limit, req.Offset)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	if len(posts) < 1 {
		return request.PostsFeedData{}, nil
	}
	files, err := p.file.GetFilesByPostUUIDs(ctx, postsIDs)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	count, err := p.post.CountRecent(ctx)
	if err != nil {
		return request.PostsFeedData{}, err
	}

	for i := range files {
		id := postIDIndexMap[files[i].ForeignUUID.String]
		posts[id].Files = append(posts[id].Files, files[i])
	}

	postDataRes := make([]request.PostDataResponse, 0, req.Limit)

	for i := range posts {
		postsFileData := make([]request.PostFileData, 0, req.Limit)
		for j := range posts[i].Files {
			postsFileData = append(postsFileData, request.PostFileData{
				Link: utils.Link(posts[i].Files[j]),
				UUID: posts[i].Files[j].UUID.String,
			})
		}

		var userData request.UserData

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

		pdr.Counts.Likes, err = p.like.CountByPost(ctx, infoblog.Like{Type: 1, ForeignUUID: types.NewNullUUID(pdr.UUID)})
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
	u := infoblog.User{UUID: types.NewNullUUID(req.UUID)}
	posts, postIDIndexMap, postsIDs, err := p.post.FindAll(ctx, u, req.Limit, req.Offset)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	if len(posts) < 1 {
		return request.PostsFeedData{}, nil
	}
	files, err := p.file.GetFilesByPostUUIDs(ctx, postsIDs)
	if err != nil {
		return request.PostsFeedData{}, err
	}
	count, err := p.post.CountByUser(ctx, u)
	if err != nil {
		return request.PostsFeedData{}, err
	}

	for i := range files {
		id := postIDIndexMap[files[i].ForeignUUID.String]
		posts[id].Files = append(posts[id].Files, files[i])
	}

	postDataRes := make([]request.PostDataResponse, 0, req.Limit)

	for i := range posts {
		postsFileData := make([]request.PostFileData, 0, req.Limit)
		for j := range posts[i].Files {
			postsFileData = append(postsFileData, request.PostFileData{
				Link: utils.Link(posts[i].Files[j]),
				UUID: posts[i].Files[j].UUID.String,
			})
		}

		var userData request.UserData

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

		pdr.Counts.Likes, err = p.like.CountByPost(ctx, infoblog.Like{Type: 1, ForeignUUID: types.NewNullUUID(pdr.UUID)})
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

func (p *Post) Like(ctx context.Context, req request.LikeReq) (request.PostDataResponse, error) {
	u, ok := ctx.Value(types.User{}).(*infoblog.User)
	if !ok {
		return request.PostDataResponse{}, fmt.Errorf("get user from request context err")
	}

	post, err := p.post.Find(ctx, infoblog.Post{
		PostEntity: infoblog.PostEntity{UUID: types.NewNullUUID(req.UUID)},
	})

	if err != nil {
		return request.PostDataResponse{}, err
	}
	like := infoblog.Like{
		UUID:        types.NewNullUUID(),
		Active:      types.NewNullBool(req.Active),
		Type:        1,
		ForeignUUID: post.UUID,
		UserUUID:    post.UserUUID,
		LikerUUID:   u.UUID,
	}

	err = p.like.Upsert(ctx, like)

	if err != nil {
		return request.PostDataResponse{}, err
	}

	var ops string

	ops = "decr"

	if req.Active {
		ops = "incr"
	}

	post, err = p.post.Count(ctx, post, "likes", ops)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	user, err := p.services.User.Count(ctx, infoblog.User{UUID: post.UserUUID}, "likes", ops)
	if err != nil {
		return request.PostDataResponse{}, err
	}
	_ = user

	postDataRes := request.PostDataResponse{}

	err = p.MapStructs(&postDataRes, &post.PostEntity)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	err = p.MapStructs(&postDataRes.User, &user)
	if err != nil {
		return request.PostDataResponse{}, err
	}

	// 4. update post and activate
	return postDataRes, nil
}

func (p *Post) Increment(ctx context.Context) (infoblog.Post, error) {
	post, err := p.post.First(ctx)
	if err != nil {
		return post, err
	}

	post, err = p.post.Count(ctx, post, "likes", "incr")
	if err != nil {
		return post, err
	}

	return post, nil
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

func (p *Post) CheckFilePermission(ctx context.Context, file infoblog.File) bool {
	var post infoblog.Post
	post.UUID = file.ForeignUUID
	post, err := p.post.Find(ctx, post)
	if err != nil {
		return false
	}

	if post.Type == PostTypeFree {
		return true
	}
	subscriber, err := extractUser(ctx)
	if err != nil || !subscriber.UUID.Valid {
		return false
	}

	if post.UserUUID.String == subscriber.UUID.String {
		return true
	}

	var user infoblog.User
	user.UUID = post.UserUUID

	return p.subscribes.CheckSubscribed(ctx, user, subscriber)
}
