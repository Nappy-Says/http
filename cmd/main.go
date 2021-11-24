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

	srv.Register(
		"/payments/{id}",
		func(req *server.Request) {
			id := req.PathParams["id"]
			log.Println(id)
		},
	)
	srv.Register(
		"/polling/{id}",
		func(req *server.Request) {
			id := req.PathParams["id"]
			log.Println(id)
		},
	)

	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// })


	return srv.Start()
}
