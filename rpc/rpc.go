package rpc

import (
	"encoding/json"
	"log"
	"net"
)

type GetFeatureRPC struct {
	RPCType string `json:"type"`
}

func setupConnection(addrAndPort string) net.Conn {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func sendRPC(conn net.Conn, rpc interface{}) interface{} {
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(rpc)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(conn)
	var res interface{}
	decoder.Decode(&res)
	return res

}

func GetFeatures(addrAndPort string) interface{} {
	conn := setupConnection(addrAndPort)
	defer conn.Close()
	return sendRPC(conn, GetFeatureRPC{"GetFeatures"})
}
