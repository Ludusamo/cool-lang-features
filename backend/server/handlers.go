package server

import (
	"cool-lang-features/rpc"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type InternalError struct {
	reason string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error: %s", e.reason)
}

/** Helper function to check for fatal errors
 * @param err the error pointer to be checked
 */
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/** Helper function attempts to encode a rpc response and checks for error
 * @param encoder encoder to write the response to
 * @param res the rpc response object
 */
func tryEncode(c net.Conn, res rpc.RPCRes) error {
	encoder := json.NewEncoder(c)
	encodingError := encoder.Encode(res)
	return encodingError
}

/** Registers the rpc handlers to the server
 * @lhs server pointer
 */
func (s *Server) RegisterHandlers() {
	s.rpcHandlers["HeartbeatSub"] = HeartbeatSubHandler
	s.rpcHandlers["GetFeatures"] = GetFeaturesHandler
	s.rpcHandlers["PostFeature"] = PostFeatureHandler
	s.rpcHandlers["DeleteFeature"] = DeleteFeatureHandler
	s.rpcHandlers["PatchFeature"] = PatchFeatureHandler
	s.rpcHandlers["GetFeature"] = GetFeatureHandler
}

/** Handler for heartbeat subscription request
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func HeartbeatSubHandler(s *Server,
	c net.Conn,
	rpcMsg rpc.RPCMapping) error {
	go Heartbeat(s, c)
	return tryEncode(c, rpc.RPCRes{nil, ""})
}

/** Handler for getting features
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func GetFeaturesHandler(s *Server,
	c net.Conn,
	rpcMsg rpc.RPCMapping) error {
	return tryEncode(c, rpc.RPCRes{s.db.GetFeatures(), ""})
}

/** Handler for posting a new feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func PostFeatureHandler(s *Server,
	c net.Conn,
	rpcMsg rpc.RPCMapping) error {
	encoder := json.NewEncoder(c)
	feat, err := s.db.AddFeature(
		rpcMsg["Name"].(string),
		rpcMsg["Description"].(string))
	if err != nil {
		encoder.Encode(rpc.RPCRes{nil, err.Error()})
		return &InternalError{"failed to add feature"}
	}
	return tryEncode(c, rpc.RPCRes{feat, ""})
}

/** Handler for deleting a feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func DeleteFeatureHandler(s *Server,
	c net.Conn,
	rpcMsg rpc.RPCMapping) error {
	s.db.DeleteFeature(int(rpcMsg["id"].(float64)))
	return tryEncode(c, rpc.RPCRes{nil, ""})
}

/** Handler for patching a feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func PatchFeatureHandler(s *Server,
	c net.Conn,
	rpcMsg rpc.RPCMapping) error {
	encoder := json.NewEncoder(c)
	feat, err := s.db.ModifyFeature(
		int(rpcMsg["id"].(float64)),
		rpcMsg["Name"].(string),
		rpcMsg["Description"].(string))
	if err != nil {
		encoder.Encode(rpc.RPCRes{nil, err.Error()})
		return &InternalError{"failed to modify feature"}
	}
	return tryEncode(c, rpc.RPCRes{feat, ""})
}

/** Handler for getting a particular feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func GetFeatureHandler(s *Server,
	c net.Conn,
	rpcMsg rpc.RPCMapping) error {
	encoder := json.NewEncoder(c)
	feat, err := s.db.GetFeature(int(rpcMsg["id"].(float64)))
	if err != nil {
		encoder.Encode(rpc.RPCRes{nil, err.Error()})
		return &InternalError{"failed to get feature"}
	}
	return tryEncode(c, rpc.RPCRes{feat, ""})
}
