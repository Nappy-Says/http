package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/Nappy-Says/http/pkg/banners"
)



type Server struct {
	mux 		*http.ServeMux
	bannersSvc	*banners.Service
}


func NewServer(mux *http.ServeMux, banners *banners.Service) *Server {
	return &Server{mux: mux, bannersSvc: banners}
}

func (s *Server) ServeHTTP(write http.ResponseWriter, request *http.Request)  {
	s.mux.ServeHTTP(write, request)
}

func (s *Server) Init()  {
	s.mux.HandleFunc("/banners.save",		s.handleSaveBanner)
	s.mux.HandleFunc("/banners.getAll", 	s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", 	s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveBannerByID)
}



func (s *Server) handleGetAllBanners(write http.ResponseWriter, request *http.Request) {
	banners, err := s.bannersSvc.All(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	write.Header().Set("Content-Type", "application/json")

	_, err = write.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(403), 403)
	}
}

func (s *Server) handleGetBannerByID(write http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	banner, err := s.bannersSvc.ByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	data, err := json.Marshal(banner)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	write.Header().Set("Contetn-Type", "applicatrion/json")

	_, err = write.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(403), 403)
	}
}

func (s *Server) handleSaveBanner(write http.ResponseWriter, request *http.Request) {
	idParam := request.PostFormValue("id")
	linkParam := request.PostFormValue("link")
	titleParam := request.PostFormValue("title")
	buttonParam := request.PostFormValue("button")
	contentParam := request.PostFormValue("content")

	log.Println("idParam: ", idParam)
	log.Println("linkParam: ", linkParam)
	log.Println("titleParam: ", titleParam)
	log.Println("buttonParam: ", buttonParam)
	log.Println("contentParam: ", contentParam)

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	banner := banners.Banner{
		ID:      id,
		Link:    linkParam,
		Title:   titleParam,
		Button:  buttonParam,
		Content: contentParam,
	}

	file, header, err := request.FormFile("image")

	log.Println(file, header)

	if err == nil {
		banner.Image = strings.Split(header.Filename, ".")[len(strings.Split(header.Filename, "."))-1]
	}

	banners, err := s.bannersSvc.Save(request.Context(), &banner, file)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	write.Header().Set("Contetn-Type", "applicatrion/json")

	_, err = write.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(403), 403)
	}
}

func (s *Server) handleRemoveBannerByID(write http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(400), 400)
		return
	}

	banners, err := s.bannersSvc.RemoveByID(request.Context(), id)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	data, err := json.Marshal(banners)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(500), 500)
		return
	}

	write.Header().Set("Contetn-Type", "applicatrion/json")

	_, err = write.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(write, http.StatusText(403), 403)
	}
}
