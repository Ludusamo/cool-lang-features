package main

import (
	"fmt"
	"github.com/Ludusamo/cool-lang-features/server"
	"github.com/pborman/getopt"
)

func main() {
	portFlag := getopt.IntLong("listen", 'p', 8080, "port number for HTTP requests")
	getopt.Parse()
	fmt.Println(*portFlag)
	svr := server.CreateServer()
	svr.RegisterHandlers()
	svr.Start(*portFlag)
}
