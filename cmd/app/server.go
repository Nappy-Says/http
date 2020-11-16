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
	log.Println("Init method")
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerById)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveById)
}

func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, err := s.bannersSvc.All(request.Context())
	requestError(writer, err, http.StatusInternalServerError)

	data, err := json.Marshal(items)
	requestError(writer, err, http.StatusInternalServerError)

	jsonResponse(writer, data)
}

func (s *Server) handleGetBannerById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	requestError(writer, err, http.StatusBadRequest)

	item, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//requestError(writer, err, http.StatusInternalServerError)

	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//requestError(writer, err, http.StatusInternalServerError)

	jsonResponse(writer, data)
}

func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	requestError(writer, err, http.StatusBadRequest)

	banner := &banners.Banner{
		ID:      id,
		Title:   request.URL.Query().Get("title"),
		Content: request.URL.Query().Get("content"),
		Button:  request.URL.Query().Get("button"),
		Link:    request.URL.Query().Get("link"),
	}

	item, err := s.bannersSvc.Save(request.Context(), banner)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	jsonResponse(writer, data)
}

func (s *Server) handleRemoveById(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	requestError(writer, err, http.StatusBadRequest)

	item, err := s.bannersSvc.RemoveByID(request.Context(), id)
	requestError(writer, err, http.StatusInternalServerError)

	data, err := json.Marshal(item)
	requestError(writer, err, http.StatusInternalServerError)

	jsonResponse(writer, data)
}

func jsonResponse(writer http.ResponseWriter, data []byte) {
	writer.Header().Set("Content-Type", "application/json")
	_, err := writer.Write(data)
	if err != nil {
		log.Println("Error write response: ", err)
	}
}

func requestError(writer http.ResponseWriter, err error, status int) {
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(status), status)
		return
	}
}
