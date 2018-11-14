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

/** Sends an rpc command and returns the response
 * @param conn the connection to send the RPC on
 * @param rpc the rpc to send over the connection
 */
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

/** Constructs an RPC to get features
 * @param addrAndPort the IP address and port to start the connection
 */
func GetFeatures(addrAndPort string) RPCRes {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		return RPCRes{nil, "Failed to open connection"}
	}
	defer conn.Close()
	return sendRPC(conn, GetFeaturesRPC{"GetFeatures"})
}

type PostFeatureRPC struct {
	RPCType     string `json:"type"`
	Name        string
	Description string
}

/** Constructs an RPC to add a new feature
 * @param addrAndPort the IP address and port to start the connection
 * @param name string identifier of the new feature
 * @param description string description of the feature
 */
func PostFeature(addrAndPort string, name string, description string) RPCRes {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		return RPCRes{nil, "Failed to open connection"}
	}
	defer conn.Close()
	return sendRPC(conn, PostFeatureRPC{"PostFeature", name, description})
}

type DeleteFeatureRPC struct {
	RPCType string `json:"type"`
	ID      int    `json:"id"`
}

/** Constructs an RPC to delete a feature
 * @param addrAndPort the IP address and port to start the connection
 * @param id int identifier of the feature
 */
func DeleteFeature(addrAndPort string, id int) RPCRes {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		return RPCRes{nil, "Failed to open connection"}
	}
	defer conn.Close()
	return sendRPC(conn, DeleteFeatureRPC{"DeleteFeature", id})
}

type PatchFeatureRPC struct {
	RPCType     string `json:"type"`
	ID          int    `json:"id"`
	Name        string
	Description string
}

/** Constructs an RPC to modify an existing feature
 * @param addrAndPort the IP address and port to start the connection
 * @param id int identifier of the feature
 * @param name new string identifier of the feature
 * @param description newstring description of the feature
 */
func PatchFeature(addrAndPort string, id int, name string, des string) RPCRes {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		return RPCRes{nil, "Failed to open connection"}
	}
	defer conn.Close()
	return sendRPC(conn, PatchFeatureRPC{"PatchFeature", id, name, des})
}

type GetFeatureRPC struct {
	RPCType string `json:"type"`
	ID      int    `json:"id"`
}

/** Constructs an RPC to retrieve a specific feature
 * @param addrAndPort the IP address and port to start the connection
 * @param id int identifier of the feature
 */
func GetFeature(addrAndPort string, id int) RPCRes {
	conn, err := net.Dial("tcp", addrAndPort)
	if err != nil {
		return RPCRes{nil, "Failed to open connection"}
	}
	defer conn.Close()
	return sendRPC(conn, GetFeatureRPC{"GetFeature", id})
}
