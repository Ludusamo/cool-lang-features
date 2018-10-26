package backend

import (
	"cool-lang-features/database"
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type Request struct {
}

type Server struct {
	db *database.Database
}

func CreateServer() *Server {
	return &Server{database.CreateDatabase()}
}

func (s *Server) Start(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		// Later on we can make this concurrent, but it can only handle one
		// connection for the time being
		if err != nil {
			log.Fatal(err)
		}
		s.HandleConnection(conn)
	}
}

func (s *Server) HandleConnection(c net.Conn) {
	defer c.Close()
	d := json.NewDecoder(c)
	var rpc map[string]interface{}
	err := d.Decode(&rpc)
	if err != nil {
		log.Fatal(err)
	}

	rpcType, exists := rpc["type"]
	if !exists {
		return
	}
	log.Println(rpcType)
	encoder := json.NewEncoder(c)
	// Handle RPC
	if rpcType == "GetFeatures" {
		err := encoder.Encode(s.db.GetFeatures())
		if err != nil {
			log.Fatal(err)
		}
	}
}

/** Adds dummy data to the database for so there are a few data points
 * @lhs server pointer
 */
func (s *Server) AddDummyData() {
	s.db.AddFeature("Pattern Matching", "Pattern matching is a tool in "+
		"programming languages to process data based on its structure.")
	s.db.AddFeature("Reflection", "Reflection is a method by which a program "+
		"can achieve metaprogramming capabilities.")
}
