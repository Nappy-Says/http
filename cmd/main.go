package main

import (
	"github.com/Nappy-Says/http/pkg/server"
	"net"
	"os"
)

func main() {
	if err := execute(server.HOST, server.PORT); err != nil {
		os.Exit(1)
	}
}

func execute(host, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	srv.Register("/", srv.RouteHandler("Welcome"))
	srv.Register("/about", srv.RouteHandler("Copyright 2020"))
	return srv.Start()
}
