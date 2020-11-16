package app

import (
	"encoding/json"
	"github.com/bdaler/http/pkg/banners"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	mux        *http.ServeMux
	bannersSvc *banners.Service
}

//NewServer construct
func NewServer(mux *http.ServeMux, bannersSvc *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: bannersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerById)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, err := s.bannersSvc.All(request.Context())
	errInternalServerError(writer, err)

	data, err := json.Marshal(items)
	errInternalServerError(writer, err)

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}

func (s *Server) handleGetBannerById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	errBadRequest(writer, err)

	item, err := s.bannersSvc.ByID(request.Context(), id)
	errInternalServerError(writer, err)

	data, err := json.Marshal(item)
	errInternalServerError(writer, err)

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	errBadRequest(writer, err)

	banner := &banners.Banner{
		ID:      id,
		Title:   request.URL.Query().Get("title"),
		Content: request.URL.Query().Get("content"),
		Button:  request.URL.Query().Get("button"),
		Link:    request.URL.Query().Get("link"),
	}

	item, err := s.bannersSvc.Save(request.Context(), banner)
	errInternalServerError(writer, err)

	data, err := json.Marshal(item)
	errInternalServerError(writer, err)

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}

func (s *Server) handleRemoveById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	errBadRequest(writer, err)

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	errInternalServerError(writer, err)

	data, err := json.Marshal(item)
	errInternalServerError(writer, err)

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}

func errBadRequest(writer http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
}

func errInternalServerError(writer http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
