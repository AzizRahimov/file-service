package main

import (
	"flag"
	"github.com/AzizRahimov/file-service/cmd/file-server/app"
	"github.com/AzizRahimov/file-service/pkg/services/files"
	"log"
	"net/http"
)

const mediaPath = "files"

func main() {
	flag.Parse()
	start()
}

func start() {

	mux := http.NewServeMux()


	fileSvc := files.NewFilesSvc(mediaPath)
	server := app.NewServer(
		mux,
		fileSvc,
		mediaPath,
	)

	server.InitRoutes()
	log.Fatal(http.ListenAndServe("0.0.0.0:9994", server))
}