package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Nappy-Says/http/pkg/banners"
)



type Server struct {
	mux 		*http.ServeMux
	bannersSvc	*banners.Service
}


func NewServer(mux *http.ServeMux, banners *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: banners}
}

func (s *Server) ServerHTTP(write http.ResponseWriter, request *http.Request)  {
	s.mux.ServeHTTP(write, request)
}

func (s *Server) Init()  {
	// s.mux.HandleFunc("/banners.save",		s.handleSaveBanner)
	// s.mux.HandleFunc("/banners.getAll", 	s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getByID", 	s.handleGetBannerByID)
	// s.mux.HandleFunc("/banners.removeByID", s.handleRemoveBannerByID)
}


func (s *Server) handleGetBannerByID(write http.ResponseWriter, request *http.Request)  {
	idParam := 	request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	s.bannersSvc.All(request.Context())
	
}
