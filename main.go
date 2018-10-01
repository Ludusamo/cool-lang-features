package main

import (
	"fmt"
	"github.com/pborman/getopt"
	"net/http"
    "github.com/Ludusamo/cool-lang-features/routes"
    "strconv"
)

func main() {
	portFlag := getopt.IntLong("listen", 'p', 8080, "port number for HTTP requests")
	getopt.Parse()
	fmt.Println(*portFlag)
	routes.RegisterHandlers()
	http.ListenAndServe(":" + strconv.Itoa(*portFlag), nil)
}
