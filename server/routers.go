package server

import (
	"github.com/go-chi/chi"
	"net/http"
)

func (serv *Server) applyRouters() {
	serv.mux.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Get("/data", serv.GetData)
			r.Post("/data", serv.PostData)
		})
		r.Handle("/*", http.FileServer(http.Dir(serv.conf.Root)))
	})
}
