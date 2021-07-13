package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
	"gitlab.com/InfoBlogFriends/server/request"

	"gitlab.com/InfoBlogFriends/server/service"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

var decoder = schema.NewDecoder()

type postsController struct {
	respond.Responder
	user   *service.User
	file   *service.File
	post   *service.Post
	logger *zap.Logger
}

func NewPostsController(responder respond.Responder, user *service.User, file *service.File, post *service.Post, logger *zap.Logger) *postsController {
	return &postsController{
		Responder: responder,
		user:      user,
		file:      file,
		post:      post,
		logger:    logger,
	}
}

func (a *postsController) Add() http.HandlerFunc {
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

		var postAddReq request.PostAddReq

		// r.PostForm is a map of our POST form values
		err = decoder.Decode(&postAddReq, r.PostForm)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		formFile := service.FormFile{
			File:       file,
			FileHeader: fHeader,
		}

		_, err = a.post.SavePost(r.Context(), formFile, postAddReq, &u)

		if err != nil {
			a.ErrorBadRequest(w, err)
			return
		}

		a.SendJSON(w, request.Response{
			Success: true,
			Msg:     "Данные профиля обновлены",
			Data:    postAddReq,
		})
	}
}
