package main

import (
	"log"
	"net"
	"os"
	"github.com/Nappy-Says/http/pkg/server"
)


func main()  {
	host := "127.0.0.1"
	port := "9999"

	if err := execute(host, port); err != nil {
		os.Exit(1)
	}
}


func execute(host string, port string) (err error) {
	srv := server.NewServer(net.JoinHostPort(host, port))
	
	srv.Register("/", func(conn net.Conn) {
		_, err := conn.Write(server.HeaderShortcut("Welcome to our web-site"))

		if err != nil {
			log.Println(err)
		}
	})

	srv.Register("/about", func(conn net.Conn) {
		_, err := conn.Write(server.HeaderShortcut("About Golang Academy"))

		if err != nil {
			log.Println(err)
		}
	})

	
	return srv.Start()
}
