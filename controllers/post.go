package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/services"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

var formDecoder = schema.NewDecoder()

type postsController struct {
	respond.Responder
	user   *services.User
	file   *services.File
	post   *services.Post
	logger *zap.Logger
}

func NewPostsController(responder respond.Responder, user *services.User, file *services.File, post *services.Post, logger *zap.Logger) *postsController {
	return &postsController{
		Responder: responder,
		user:      user,
		file:      file,
		post:      post,
		logger:    logger,
	}
}

func (a *postsController) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := extractUser(r)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}
		err = r.ParseMultipartForm(100 << 20)
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

		var postAddReq request.PostCreateReq

		// r.PostForm is a map of our POST form values
		err = formDecoder.Decode(&postAddReq, r.PostForm)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		formFile := services.FormFile{
			File:       file,
			FileHeader: fHeader,
		}

		post, err := a.post.SavePost(r.Context(), formFile, postAddReq, &u)

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

func (a *postsController) FeedRecent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := extractUser(r)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		var postsListReq request.PostsFeedReq
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
			Msg:     "",
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

		var postsListReq request.PostsFeedReq
		err = Decode(r, &postsListReq)
		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		feed, err := a.post.FeedMy(r.Context(), u, postsListReq)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}

		a.SendJSON(w, request.PostsFeedResponse{
			Success: true,
			Msg:     "",
			Data:    feed,
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
