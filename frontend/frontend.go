package main

import (
	"cool-lang-features/frontend/server"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func extractHostAndPort(hostAndPort string) (string, int) {
	split := strings.Split(hostAndPort, ":")
	host := split[0]
	port := 8090
	if len(split) > 1 {
		p, err := strconv.Atoi(split[1])
		if err != nil {
			log.Fatal(err)
		}
		port = p
	}
	return host, port
}

func main() {
	portFlag := flag.Int("listen", 8080, "port number for HTTP requests")
	backendList := flag.String("backend", ":8090", "hostnames and ports of backends")
	flag.Parse()
	fmt.Println(*portFlag)

	var backends []string
	for _, backend := range strings.Split(*backendList, ",") {
		host, port := extractHostAndPort(backend)
		backends = append(backends, host+":"+strconv.Itoa(port))
	}
	svr := server.CreateServer(backends)
	svr.RegisterHandlers()
	svr.Start(*portFlag)
}
