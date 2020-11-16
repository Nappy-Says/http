package main

import (
	"github.com/Nappy-Says/http/cmd/app"
	"github.com/Nappy-Says/http/pkg/banners"
	"net"
	"net/http"
	"os")

const (
	HOST = "0.0.0.0"
	PORT = "9999")

func main() {
	if err := execute(HOST, PORT); err != nil {
		os.Exit(1)
	}
}

func execute(server, port string) (err error) {
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	serverHandler := app.NewServer(mux, bannersSvc)
	serverHandler.Init()

	srv := &http.Server{
		Addr:    net.JoinHostPort(server, port),
		Handler: serverHandler,
	}
	return srv.ListenAndServe()
}
