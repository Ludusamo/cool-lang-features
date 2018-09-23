package main

import (
    "fmt"
    "github.com/pborman/getopt"
)

func main() {
    portFlag := getopt.Int32Long("listen", 'p', 8080, "port number for HTTP requests")
    getopt.Parse()
    fmt.Println(*portFlag)
}
