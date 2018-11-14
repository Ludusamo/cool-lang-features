package server

import (
	"cool-lang-features/database"
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type Server struct {
	db *database.Database
	// Mapping of handlers that takes a server, a net connection, and an rpc
	rpcHandlers    map[string]func(*Server, *json.Encoder, RPCMapping) error
	connectedUsers map[string]net.Conn
}

/** Creates an empty server object
 * @return pointer to created server
 */
func CreateServer() *Server {
	return &Server{
		database.CreateDatabase(),
		make(map[string]func(*Server, *json.Encoder, RPCMapping) error),
		make(map[string]net.Conn)}
}

/** Spins up the service to listen to external tcp requests
 * @lhs server pointer
 * @param port integer port where the service will listen from
 */
func (s *Server) Start(port int) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go s.HandleConnection(conn)
	}
}

/** Takes a connection, read an RPC, and processes it. Closes connection when
 * finished
 * @param c open connection
 */
func (s *Server) HandleConnection(c net.Conn) {
	defer c.Close()
	ip := c.RemoteAddr().String()
	log.Println(ip)
	s.connectedUsers[ip] = c
	d := json.NewDecoder(c)
	for {
		var rpcMsg RPCMapping
		d.Decode(&rpcMsg)

		rpcType, exists := rpcMsg["type"]
		if !exists {
			return
		}
		log.Println(rpcMsg)
		encoder := json.NewEncoder(c)
		// Handle RPC
		err := s.rpcHandlers[rpcType.(string)](s, encoder, rpcMsg)
		if err != nil {
			switch err.(type) {
			case *InternalError:
				break
			default:
				log.Println(err)
			}
		}
	}
	log.Printf("removing %s\n", ip)
	delete(s.connectedUsers, ip)
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
