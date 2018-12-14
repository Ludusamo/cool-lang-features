package inc

import (
	"encoding/json"
	"log"
	"net"
	"strconv"
	"sync"
)

type Transaction struct {
	TransactID int
	Voted      map[string]int
}

func CreateTransaction(id int) Transaction {
	return Transaction{id, make(map[string]int)}
}

type MsgType int

type Msg struct {
	Type    MsgType
	Payload Transaction
	PID     int
	RPID    int
	Addr    string
}

type OutgoingMessage struct {
	Msg  Msg
	Addr string
	PID  int
}

type IncomingMessage struct {
	Msg  Msg
	Addr string
	PID  int
}

type Conn struct {
	pid int
	In  chan IncomingMessage
	Out chan OutgoingMessage
	up  bool
}

func (c *Conn) Send(addr string, pid int, msg Msg) {
	c.Out <- OutgoingMessage{msg, addr, pid}
}

type ConnectionManager struct {
	conns         map[string]net.Conn
	localConns    map[int]*Conn
	connsLock     *sync.RWMutex
	newConnection chan *Conn
}

func sendQueued(manager *ConnectionManager, c *Conn) {
	for {
		outgoingMsg := <-c.Out
		log.Printf("Sending msg to %s\n", outgoingMsg.Addr)
		conn, up := manager.conns[outgoingMsg.Addr]
		if !up {
			c.up = false
			continue
		}
		outgoingMsg.Msg.Addr = conn.LocalAddr().String()
		encoder := json.NewEncoder(conn)
		err := encoder.Encode(outgoingMsg.Msg)
		if err != nil {
			c.up = false
		}
	}
}

func CreateConnectionManager(port int) *ConnectionManager {
	manager := ConnectionManager{make(map[string]net.Conn),
		make(map[int]*Conn),
		&sync.RWMutex{},
		make(chan *Conn)}
	go func() {
		ln, _ := net.Listen("tcp", ":"+strconv.Itoa(port))
		for {
			conn, err := ln.Accept()
			if err == nil {
				go handleTCPConn(&manager, conn)
			}
		}
	}()
	return &manager
}

func handleTCPConn(cm *ConnectionManager, conn net.Conn) {
	d := json.NewDecoder(conn)
	for {
		var msg Msg
		d.Decode(&msg)
		log.Println(msg.Addr)
		localConn, localConnExists := cm.localConns[msg.PID]
		if !localConnExists {
			continue
		}
		localConn.In <- IncomingMessage{msg,
			msg.Addr,
			msg.RPID}
	}
}

func (c *ConnectionManager) GetLocalConn(pid int) *Conn {
	conn, open := c.localConns[pid]
	if !open {
		c.localConns[pid] = &Conn{
			pid,
			make(chan IncomingMessage),
			make(chan OutgoingMessage),
			true}
		go sendQueued(c, c.localConns[pid])
		conn = c.localConns[pid]
	}
	return conn
}

func (c *ConnectionManager) Open(addr string) error {
	c.connsLock.RLock()
	_, open := c.conns[addr]
	c.connsLock.RUnlock()
	if !open {
		newConn, err := net.Dial("tcp", addr)
		if err != nil {
			return err
		}
		go handleTCPConn(c, newConn)
		c.connsLock.Lock()
		c.conns[addr] = newConn
		c.connsLock.Unlock()
	}
	return nil
}
