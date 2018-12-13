package server

import (
	"cool-lang-features/database"
	"cool-lang-features/inc"
	"cool-lang-features/rpc"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	db *database.Database
	// Mapping of handlers that takes a server, a net connection, and an rpc
	rpcHandlers       map[string]func(*Server, net.Conn, rpc.RPCMapping) error
	connectedUserLock *sync.RWMutex
	connectedUsers    map[string]net.Conn
	connManager       *inc.ConnectionManager
}

/** Creates an empty server object
 * @return pointer to created server
 */
func CreateServer() *Server {
	return &Server{
		database.CreateDatabase(),
		make(map[string]func(*Server, net.Conn, rpc.RPCMapping) error),
		&sync.RWMutex{},
		make(map[string]net.Conn),
		nil}
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
	d := json.NewDecoder(c)
	s.connectedUserLock.Lock()
	s.connectedUsers[ip] = c
	s.connectedUserLock.Unlock()
	for {
		var rpcMsg rpc.RPCMapping
		d.Decode(&rpcMsg)

		// Checks to see if the heartbeat sender terminated the connection
		s.connectedUserLock.RLock()
		_, stillAlive := s.connectedUsers[ip]
		s.connectedUserLock.RUnlock()
		if !stillAlive {
			break
		}

		rpcType, exists := rpcMsg["type"]
		if !exists {
			break
		}
		log.Println(rpcMsg)
		// Handle RPC
		err := s.rpcHandlers[rpcType.(string)](s, c, rpcMsg)
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
func Heartbeat(s *Server, c net.Conn) {
	ip := c.RemoteAddr().String()
	seq := 0
	for {
		c.SetDeadline(time.Now().Add(time.Second * 20))
		encoder := json.NewEncoder(c)
		err := encoder.Encode(rpc.HeartbeatRPC{"Heartbeat", seq})
		if err != nil {
			log.Printf("removing %s\n", ip)
			s.connectedUserLock.Lock()
			delete(s.connectedUsers, ip)
			s.connectedUserLock.Unlock()
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
