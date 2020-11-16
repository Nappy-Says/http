package main

import (
	"github.com/bdaler/http/cmd/app"
	"github.com/bdaler/http/pkg/banners"
	"github.com/bdaler/http/pkg/server"
	"net"
	"net/http"
	"os"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

func execute() (err error) {
	mux := http.NewServeMux()
	bannersSvc := banners.NewService()
	serverHandler := app.NewServer(mux, bannersSvc)

	srv := &http.Server{
		Addr:    net.JoinHostPort(server.HOST, server.PORT),
		Handler: serverHandler,
	}
	return srv.ListenAndServe()

}
