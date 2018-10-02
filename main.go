package main

import (
	"fmt"
	"github.com/Ludusamo/cool-lang-features/routes"
	"github.com/Ludusamo/cool-lang-features/server"
	"github.com/pborman/getopt"
	"net/http"
	"strconv"
)

func main() {
	portFlag := getopt.IntLong("listen", 'p', 8080, "port number for HTTP requests")
	getopt.Parse()
	fmt.Println(*portFlag)
	server := CreateServer()
	server.RegisterHandlers()
	http.ListenAndServe(":"+strconv.Itoa(*portFlag), server.router)
}
