package controllers

import (
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"

	"gitlab.com/InfoBlogFriends/server/services"

	"github.com/go-chi/chi/v5"

	"gitlab.com/InfoBlogFriends/server/respond"
	"go.uber.org/zap"
)

type fileController struct {
	respond.Responder
	file   *services.File
	post   *services.Post
	logger *zap.Logger
}

func NewFileController(responder respond.Responder, services *services.Services, logger *zap.Logger) *fileController {
	return &fileController{
		Responder: responder,
		file:      services.File,
		post:      services.Post,
		logger:    logger,
	}
}

func (a *fileController) GetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileID := chi.URLParam(r, "fileID")
		file, err := a.file.GetFile(r.Context(), fileID)
		if err != nil {
			a.ErrorInternal(w, err)
			return
		}
		path := strings.Join([]string{".", file.Dir, file.Name}, "/")
		fileRaw, err := os.Open(path)
		defer func() {
			err = fileRaw.Close() //Close after function return
			if err != nil {
				log.Printf("Ошибка при: %v\n", err)
			}
		}()
		if err != nil {
			//File not found, send 404
			w.WriteHeader(http.StatusNotFound)
			_, err = w.Write([]byte("Файл не найден на сервере"))
			if err != nil {
				log.Printf("Ошибка при: %v\n", err)
			}

			return
		}

		//Create a buffer to store the header of the file in
		FileHeader := make([]byte, 512)
		//Copy the headers into the FileHeader buffer
		_, err = fileRaw.Read(FileHeader)
		if err != nil {
			log.Printf("Ошибка при: %v\n", err)
		}
		//Get content type of file
		FileContentType := http.DetectContentType(FileHeader)
		//Get the file size
		FileStat, _ := fileRaw.Stat()
		FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get
		//Send the headers

		if file.Type == services.FileTypePost {
			if a.post.CheckFilePermission(r.Context(), file) {
				_, err = fileRaw.Seek(0, 0)
				w.Header().Set("Content-Disposition", "inline;")
				w.Header().Set("Content-Type", FileContentType)
				w.Header().Set("Content-Length", FileSize)
				_, _ = io.Copy(w, fileRaw)
				return
			}

			if FileContentType == "image/jpeg" {
				_, err = fileRaw.Seek(0, 0)
				if err != nil {
					a.ErrorInternal(w, err)
					return
				}
				m, _, err := image.Decode(fileRaw)
				if err != nil {
					a.ErrorInternal(w, err)
					return
				}
				dstImage := imaging.Blur(m, 21)
				if err := jpeg.Encode(w, dstImage, nil); err != nil {
					log.Printf("failed to encode: %v", err)
				}

				return
			}
			//We read 512 bytes from the file already, so we reset the offset back to 0
			_, err = fileRaw.Seek(0, 0)
			w.Header().Set("Content-Disposition", "inline;")
			w.Header().Set("Content-Type", FileContentType)
			w.Header().Set("Content-Length", FileSize)
			_, _ = io.Copy(w, fileRaw)
			return
		}

		_, _ = io.Copy(w, fileRaw)
		return
	}
}
