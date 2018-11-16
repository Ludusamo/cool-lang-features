package rpc

import (
	"encoding/json"
	"log"
	"net"
)

type RPCMapping map[string]interface{}

type RPCRes struct {
	Data interface{}
	Err  string
}

/** Sends an rpc command and returns the response
 * @param conn the connection to send the RPC on
 * @param rpc the rpc to send over the connection
 */
func SendRPC(conn net.Conn, rpc interface{}) RPCRes {
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

type HeartbeatSubRPC struct {
	RPCType string `json:"type"`
}

type HeartbeatRPC struct {
	RPCType string `json:"type"`
	SeqNum  int    `json:"seq"`
}

type GetFeaturesRPC struct {
	RPCType string `json:"type"`
}

type PostFeatureRPC struct {
	RPCType     string `json:"type"`
	Name        string
	Description string
}

type DeleteFeatureRPC struct {
	RPCType string `json:"type"`
	ID      int    `json:"id"`
}

type PatchFeatureRPC struct {
	RPCType     string `json:"type"`
	ID          int    `json:"id"`
	Name        string
	Description string
}

type GetFeatureRPC struct {
	RPCType string `json:"type"`
	ID      int    `json:"id"`
}
