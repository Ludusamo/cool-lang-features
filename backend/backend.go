package main

import (
	"cool-lang-features/backend/server"
	"flag"
	"fmt"
)

func main() {
	portFlag := flag.Int("listen", 8090, "port number for HTTP requests")
	flag.Parse()
	fmt.Println(*portFlag)
	svr := server.CreateServer()
	svr.AddDummyData()
	svr.Start(*portFlag)
}
