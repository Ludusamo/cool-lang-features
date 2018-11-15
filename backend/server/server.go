package server

import (
	"cool-lang-features/database"
	"cool-lang-features/rpc"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"
)

type Server struct {
	db *database.Database
	// Mapping of handlers that takes a server, a net connection, and an rpc
	rpcHandlers    map[string]func(*Server, *json.Encoder, rpc.RPCMapping) error
	connectedUsers map[string]net.Conn
}

/** Creates an empty server object
 * @return pointer to created server
 */
func CreateServer() *Server {
	return &Server{
		database.CreateDatabase(),
		make(map[string]func(*Server, *json.Encoder, rpc.RPCMapping) error),
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
	go heartbeat(s, c)
	for {
		var rpcMsg rpc.RPCMapping
		d.Decode(&rpcMsg)

		// Checks to see if the heartbeat sender terminated the connection
		_, stillAlive := s.connectedUsers[ip]
		if !stillAlive {
			break
		}

		rpcType, exists := rpcMsg["type"]
		if !exists {
			break
		}
		log.Println(rpcMsg)
		encoder := json.NewEncoder(c)
		// Handle RPC
		err := s.rpcHandlers[rpcType.(string)](s, encoder, rpcMsg)
		if err != nil {
			switch err.(type) {
			case *InternalError:
				log.Println(err)
			default:
				continue
			}
		}
	}
}

/** Sends frontend a heartbeat rpc to let them know that the server is up
 * @param s server pointer
 * @param c connection to send the heartbeat across
 */
func heartbeat(s *Server, c net.Conn) {
	seq := 0
	for {
		c.SetDeadline(time.Now().Add(time.Second * 20))
		encoder := json.NewEncoder(c)
		err := encoder.Encode(rpc.HeartbeatRPC{"Heartbeat", seq})
		if err != nil {
			ip := c.RemoteAddr().String()
			log.Printf("removing %s\n", ip)
			delete(s.connectedUsers, ip)
			return
		}
		seq += 1
		time.Sleep(time.Second * 10)
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
