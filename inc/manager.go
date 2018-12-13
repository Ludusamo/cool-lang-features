package inc

import (
	"cool-lang-features/rpc"
	"encoding/json"
	"net"
	"strconv"
	"sync"
)

type Conn struct {
	pid  int
	in   chan rpc.RPCMapping
	out  chan rpc.RPCMapping
	addr string
	up   bool
}

type ConnectionManager struct {
	conns      map[string]net.Conn
	localConns map[int]*Conn
	connsLock  *sync.RWMutex
}

func sendQueued(manager *ConnectionManager, c *Conn) {
	for {
		conn, up := manager.conns[c.addr]
		if !up {
			c.up = false
			continue
		}
		msg, ok := <-c.out
		if ok {
			encoder := json.NewEncoder(conn)
			err := encoder.Encode(msg)
			if err != nil {
				c.up = false
			}
		}
	}
}

func CreateConnectionManager(port int) *ConnectionManager {
	manager := ConnectionManager{make(map[string]net.Conn),
		make(map[int]*Conn),
		&sync.RWMutex{}}
	go func() {
		ln, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
		for {
			conn, err := ln.Accept()
			if err == nil {
				addr := conn.RemoteAddr().String()
				if _, exists := manager.conns[addr]; !exists {
					manager.connsLock.Lock()
					manager.conns[addr] = conn
					manager.connsLock.Unlock()
					go handleTCPConn(&manager, conn)
				}
			}
		}
	}()
	return &manager
}

func handleTCPConn(cm *ConnectionManager, conn net.Conn) {
	d := json.NewDecoder(conn)
	for {
		var rpcMsg rpc.RPCMapping
		d.Decode(&rpcMsg)
		pid, hasId := rpcMsg["pid"]
		if !hasId {
			continue
		}
		localConn, localConnExists := cm.localConns[int(pid.(float64))]
		if !localConnExists {
			continue
		}
		localConn.in <- rpcMsg
	}
}

func (c *ConnectionManager) Open(addr string, pid int) (*Conn, error) {
	c.connsLock.RLock()
	conn, open := c.conns[addr]
	c.connsLock.RUnlock()
	if !open {
		newConn, err := net.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		go handleTCPConn(c, conn)
		c.connsLock.Lock()
		c.conns[addr] = newConn
		c.connsLock.Unlock()
		conn = newConn
	}
	c.localConns[pid] = &Conn{
		pid,
		make(chan rpc.RPCMapping),
		make(chan rpc.RPCMapping),
		addr,
		true}
	go sendQueued(c, c.localConns[pid])
	return c.localConns[pid], nil
}
