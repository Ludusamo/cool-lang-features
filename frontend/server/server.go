package server

import (
	"cool-lang-features/database"
	"cool-lang-features/rpc"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type RegexHandlerMapping struct {
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type Router struct {
	routes []RegexHandlerMapping
}

type Server struct {
	db                  *database.Database
	router              *Router
	backends            []string
	backendConnection   net.Conn
	heartbeatConnection net.Conn
	backendUp           bool
}

/** Creates an empty router object
 * @return pointer to created router
 */
func CreateRouter() *Router {
	return &Router{make([]RegexHandlerMapping, 0)}
}

/** Creates an empty server object
 * @return pointer to created server
 */
func CreateServer(backends []string) *Server {
	server := &Server{database.CreateDatabase(),
		CreateRouter(),
		backends,
		nil,
		nil,
		false}
	server.connectToBackend()
	go heartbeatMonitor(server)
	return server
}

/** Attempts to connect the server to the specified backend server
 * lhs server pointer
 */
func (s *Server) connectToBackend() {
	conn, err := net.Dial("tcp", s.backends[0])
	s.backendConnection = conn
	s.backendUp = err == nil
}

func (s *Server) startHeartbeatMonitor() {
	conn, err := net.Dial("tcp", s.backends[0])
	s.heartbeatConnection = conn
	s.backendUp = err == nil
	if s.backendUp {
		res := rpc.SendRPC(conn, rpc.HeartbeatSubRPC{"HeartbeatSub"})
		for res.Err != "" {
			res = rpc.SendRPC(conn, rpc.HeartbeatSubRPC{"HeartbeatSub"})
		}
	}
}

/** Monitors heartbeat messages from the backend; if none is received, the
 * backend is assumed to be down and the monitor will attempt to reconnect
 * @param s server pointer
 */
func heartbeatMonitor(s *Server) {
	s.startHeartbeatMonitor()
	lastTimeReceived := time.Now()
	seq := -1
	for {
		d := json.NewDecoder(s.heartbeatConnection)
		if !s.backendUp {
			fmt.Printf("Detected failure on %s at %v\n",
				s.backends[0],
				time.Now().UTC())
			s.connectToBackend()
			s.startHeartbeatMonitor()
			continue
		}

		var rpcMsg rpc.RPCMapping
		err := d.Decode(&rpcMsg)
		if err != nil {
			if time.Now().Sub(lastTimeReceived) > time.Second*35 {
				fmt.Println(err)
				s.backendConnection.Close()
				s.heartbeatConnection.Close()
				s.backendUp = false
			}
			continue
		}
		rpcType, exists := rpcMsg["type"]
		if exists && rpcType.(string) == "Heartbeat" {
			receivedSeq := int(rpcMsg["seq"].(float64))
			if receivedSeq > seq {
				seq = receivedSeq
				lastTimeReceived = time.Now()
			}
		}

		fmt.Println("received heartbeat")
	}
}

/** Spins up the service to listen to external http requests
 * @lhs server pointer
 * @param port integer port where the service will listen from
 * @param backendHostname hostname of the backend
 * @param backendPort port of the backend
 */
func (s *Server) Start(port int) {
	http.ListenAndServe(":"+strconv.Itoa(port), s.router)
}
