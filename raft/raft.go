package raft

import (
	"cool-lang-features/inc"
	"log"
	"math/rand"
	"time"
)

type State int

const (
	FOLLOWER State = iota
	CANDIDATE
	LEADER
)

const (
	RequestToVote = iota
	Vote
)

type RaftNode struct {
	connManager *inc.ConnectionManager
	backends    []string
	leader      string
	log         []inc.Transaction
	uncommitted map[int]inc.Transaction

	raftConn *inc.Conn
	state    State
}

func keepConnUp(node *RaftNode, addr string) {
	err := node.connManager.Open(addr)
	for err != nil {
		err = node.connManager.Open(addr)
	}
}

func CreateRaftNode(port int, backends []string) *RaftNode {
	node := RaftNode{inc.CreateConnectionManager(port),
		backends,
		"",
		make([]inc.Transaction, 0),
		make(map[int]inc.Transaction),
		nil,
		FOLLOWER}
	for _, backend := range backends {
		go keepConnUp(&node, backend)
	}
	node.raftConn = node.connManager.GetLocalConn(0)
	go RaftLoop(&node)
	return &node
}

func genTransactID(node *RaftNode) int {
	for {
		id := rand.Int()
		if _, exists := node.uncommitted[id]; !exists {
			return id
		}
	}
}

func RunElection(node *RaftNode) {
	node.state = CANDIDATE
	id := genTransactID(node)
	node.uncommitted[id] = inc.CreateTransaction(id)
	for _, addr := range node.backends {
		node.raftConn.Send(addr, 0, inc.Msg{RequestToVote, node.uncommitted[id], 0, 0, ""})
	}
}

func RaftLoop(node *RaftNode) {
	for {
		timeout := time.NewTimer(time.Duration(10+rand.Intn(10)) * time.Second)
		select {
		case <-timeout.C:
			log.Println("Timed Out")
			RunElection(node)
		case incoming := <-node.raftConn.In:
			msg := incoming.Msg
			if node.state == CANDIDATE {
				if msg.Type == Vote {
					log.Println("Received Vote")
				}
			} else if node.state == FOLLOWER {
				if msg.Type == RequestToVote {
					log.Println("Received Request to Vote")
					node.raftConn.Send(incoming.Addr, incoming.PID, inc.Msg{Vote, msg.Payload, 0, 0, ""})
				}
			}
		}
	}
}
