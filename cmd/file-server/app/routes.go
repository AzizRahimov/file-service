package app

import "net/http"

func (s *server) InitRoutes() {
	mux := s.router.(*http.ServeMux)

	mux.HandleFunc("/api/files", s.handleMultipart)
	mux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(s.storagePath))))


}
