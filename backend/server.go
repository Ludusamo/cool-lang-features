package backend

import (
	"cool-lang-features/database"
	"cool-lang-features/rpc"
	"encoding/json"
	"log"
	"net"
	"strconv"
)

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
	var rpcMsg map[string]interface{}
	err := d.Decode(&rpcMsg)
	if err != nil {
		log.Fatal(err)
	}

	rpcType, exists := rpcMsg["type"]
	if !exists {
		return
	}
	log.Println(rpcMsg)
	encoder := json.NewEncoder(c)
	// Handle RPC
	if rpcType == "GetFeatures" {
		err := encoder.Encode(rpc.RPCRes{s.db.GetFeatures(), ""})
		if err != nil {
			log.Fatal(err)
		}
	} else if rpcType == "PostFeature" {
		feat, err := s.db.AddFeature(
			rpcMsg["Name"].(string),
			rpcMsg["Description"].(string))
		if err != nil {
			encoder.Encode(rpc.RPCRes{nil, err.Error()})
			return
		}
		encodingError := encoder.Encode(rpc.RPCRes{feat, ""})
		if encodingError != nil {
			log.Fatal(encodingError)
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
