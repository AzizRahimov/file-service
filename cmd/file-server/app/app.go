package app

import (
	"encoding/json"
	"errors"
	"github.com/AzizRahimov/file-service/pkg/services/files"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type server struct {
	router        http.Handler
	fileSvc       *files.FileService
	storagePath   string
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func NewServer(router http.Handler, fileSvc *files.FileService, storagePath string) *server {
	if router == nil {
		panic(errors.New("router can't be nil"))
	}
	if fileSvc == nil {
		panic(errors.New("fileSvc can't be nil"))
	}
	if storagePath == "" {
		panic(errors.New("storagePath can't be nil"))
	}

	return &server{fileSvc: fileSvc, storagePath: storagePath, router: router}
}

const multipartMaxBytes = 10 * 1024 * 1024

func (s *server) handleMultipart(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		return
	}

	err := request.ParseMultipartForm(multipartMaxBytes)
	if err != nil {
		log.Print(err)
		http.Error(responseWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fileHeaders := request.MultipartForm.File["file"]

	type FileURL struct {
		Id   string `json:"id"`
		Path string `json:"path"`
	}
	fileURLs := make([]FileURL, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		name, err := s.saveFile(fileHeader)
		if err != nil {
			http.Error(responseWriter, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			log.Print(err)
			return
		}

		fileURLs = append(fileURLs, FileURL{
			Id:   name[:len(name)-len(filepath.Ext(name))],
			Path: "/" + s.storagePath + "/" + name,
		})
	}

	urlsJSON, err := json.Marshal(fileURLs)
	if err != nil {
		http.Error(responseWriter, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	_, err = responseWriter.Write(urlsJSON)
	if err != nil {
		log.Print(err)
		return
	}

	return
}

func (s *server) saveFile(fileHeader *multipart.FileHeader) (name string, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return
	}
	defer func() {
		err = file.Close()
	}()
// content-type - какой формат отправляется

	contentType := fileHeader.Header.Get("Content-Type")
	name, err = s.fileSvc.Save(file, contentType)
	if err != nil {
		return
	}

	return
}