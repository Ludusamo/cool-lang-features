package main

import (
	"cool-lang-features/server"
	"flag"
	"fmt"
)

func main() {
	portFlag := flag.Int("listen", 8080, "port number for HTTP requests")
	flag.Parse()
	fmt.Println(*portFlag)
	svr := server.CreateServer()
	svr.RegisterHandlers()
	svr.AddDummyData()
	svr.Start(*portFlag)
}
