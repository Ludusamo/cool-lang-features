package rpc

import (
	"encoding/json"
	"log"
	"net"
)

type RPCRes struct {
	Data interface{}
	Err  string
}

func setupConnection(addrAndPort string) net.Conn {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func sendRPC(conn net.Conn, rpc interface{}) RPCRes {
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(rpc)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(conn)
	var res RPCRes
	decoder.Decode(&res)
	return res

}

type GetFeaturesRPC struct {
	RPCType string `json:"type"`
}

func GetFeatures(addrAndPort string) RPCRes {
	conn := setupConnection(addrAndPort)
	defer conn.Close()
	return sendRPC(conn, GetFeaturesRPC{"GetFeatures"})
}

type PostFeatureRPC struct {
	RPCType     string `json:"type"`
	Name        string
	Description string
}

func PostFeature(addrAndPort string, name string, description string) RPCRes {
	conn := setupConnection(addrAndPort)
	defer conn.Close()
	return sendRPC(conn, PostFeatureRPC{"PostFeature", name, description})
}

type DeleteFeatureRPC struct {
	RPCType string `json:"type"`
	ID      int    `json:"id"`
}

func DeleteFeature(addrAndPort string, id int) RPCRes {
	conn := setupConnection(addrAndPort)
	defer conn.Close()
	return sendRPC(conn, DeleteFeatureRPC{"DeleteFeature", id})
}

type PatchFeatureRPC struct {
	RPCType     string `json:"type"`
	ID          int    `json:"id"`
	Name        string
	Description string
}

func PatchFeature(addrAndPort string, id int, name string, des string) RPCRes {
	conn := setupConnection(addrAndPort)
	defer conn.Close()
	return sendRPC(conn, PatchFeatureRPC{"PatchFeature", id, name, des})
}

type GetFeatureRPC struct {
	RPCType string `json:"type"`
	ID      int    `json:"id"`
}

func GetFeature(addrAndPort string, id int) RPCRes {
	conn := setupConnection(addrAndPort)
	defer conn.Close()
	return sendRPC(conn, GetFeatureRPC{"GetFeature", id})
}
