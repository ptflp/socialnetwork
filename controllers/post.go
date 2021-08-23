package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type postsController struct {
	*decoder.Decoder
	respond.Responder
	user   *services.User
	file   *services.File
	post   *services.Post
	logger *zap.Logger
}

func NewPostsController(responder respond.Responder, user *services.User, file *services.File, post *services.Post, logger *zap.Logger) *postsController {
	return &postsController{
		Decoder:   decoder.NewDecoder(),
		Responder: responder,
		user:      user,
		file:      file,
		post:      post,
		logger:    logger,
	}
}

func (a *postsController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postAddReq request.PostCreateReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		post, err := a.post.SavePost(r.Context(), postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    post,
		})
	}
}

func (a *postsController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postAddReq request.PostUpdateReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.post.Update(r.Context(), postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "пост успешно обновлен",
		})
	}
}

func (a *postsController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postAddReq request.PostUUIDReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.post.Delete(r.Context(), postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "пост успешно удален",
		})
	}
}

func (a *postsController) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postAddReq request.PostUUIDReq

		// r.PostForm is a map of our POST form values
		err := Decode(r, &postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		post, err := a.post.Get(r.Context(), postAddReq)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    post,
		})
	}
}

func (a *postsController) UploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(100 << 20)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		file, fHeader, err := r.FormFile("file")
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		defer file.Close()

		formFile := services.FormFile{
			File:       file,
			FileHeader: fHeader,
		}

		fileData, err := a.post.SaveFile(r.Context(), formFile)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    fileData,
		})
	}
}

func (a *postsController) FeedRecent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var postsListReq request.LimitOffsetReq
		err = Decode(r, &postsListReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		feed, err := a.post.FeedRecent(r.Context(), postsListReq)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.PostsFeedResponse{
			Success: true,
			Data:    feed,
		})
	}
}

func (a *postsController) FeedMy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := extractUser(r)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		req := request.PostsFeedUserReq{}
		err = Decode(r, &req)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		req.UUID = u.UUID.String

		feed, err := a.post.FeedByUser(r.Context(), req)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.PostsFeedResponse{
			Success: true,
			Data:    feed,
		})
	}
}

func (a *postsController) FeedUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postsListReq request.PostsFeedUserReq
		err := Decode(r, &postsListReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		feed, err := a.post.FeedByUser(r.Context(), postsListReq)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.PostsFeedResponse{
			Success: true,
			Data:    feed,
		})
	}
}

func (a *postsController) TestIncrement() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		feed, err := a.post.Increment(r.Context())
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Data:    feed,
		})
	}
}

func (a *postsController) Like() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var postLikeReq request.PostUUIDReq
		err := Decode(r, &postLikeReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		err = a.post.Like(r.Context(), postLikeReq)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.PostsFeedResponse{
			Success: true,
		})
	}
}

func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
