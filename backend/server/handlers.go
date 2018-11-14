package server

import (
	"cool-lang-features/rpc"
	"encoding/json"
	"log"
)

type RPCMapping map[string]interface{}

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
func tryEncode(encoder *json.Encoder, res rpc.RPCRes) {
	encodingError := encoder.Encode(res)
	check(encodingError)
}

/** Registers the rpc handlers to the server
 * @lhs server pointer
 */
func (s *Server) RegisterHandlers() {
	s.rpcHandlers["GetFeatures"] = GetFeaturesHandler
	s.rpcHandlers["PostFeature"] = PostFeatureHandler
	s.rpcHandlers["DeleteFeature"] = DeleteFeatureHandler
	s.rpcHandlers["PatchFeature"] = PatchFeatureHandler
	s.rpcHandlers["GetFeature"] = GetFeatureHandler
}

/** Handler for getting features
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func GetFeaturesHandler(s *Server, encoder *json.Encoder, rpcMsg RPCMapping) {
	tryEncode(encoder, rpc.RPCRes{s.db.GetFeatures(), ""})
}

/** Handler for posting a new feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func PostFeatureHandler(s *Server, encoder *json.Encoder, rpcMsg RPCMapping) {
	feat, err := s.db.AddFeature(
		rpcMsg["Name"].(string),
		rpcMsg["Description"].(string))
	if err != nil {
		encoder.Encode(rpc.RPCRes{nil, err.Error()})
		return
	}
	tryEncode(encoder, rpc.RPCRes{feat, ""})
}

/** Handler for deleting a feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func DeleteFeatureHandler(s *Server, encoder *json.Encoder, rpcMsg RPCMapping) {
	s.db.DeleteFeature(int(rpcMsg["id"].(float64)))
	tryEncode(encoder, rpc.RPCRes{nil, ""})
}

/** Handler for patching a feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func PatchFeatureHandler(s *Server, encoder *json.Encoder, rpcMsg RPCMapping) {
	feat, err := s.db.ModifyFeature(
		int(rpcMsg["id"].(float64)),
		rpcMsg["Name"].(string),
		rpcMsg["Description"].(string))
	if err != nil {
		encoder.Encode(rpc.RPCRes{nil, err.Error()})
		return
	}
	tryEncode(encoder, rpc.RPCRes{feat, ""})
}

/** Handler for getting a particular feature
 * @param s server pointer
 * @param encoder the encoder to write resposnes
 * @param rpcMsg rpc to be handled
 */
func GetFeatureHandler(s *Server, encoder *json.Encoder, rpcMsg RPCMapping) {
	feat, err := s.db.GetFeature(int(rpcMsg["id"].(float64)))
	if err != nil {
		encoder.Encode(rpc.RPCRes{nil, err.Error()})
		return
	}
	tryEncode(encoder, rpc.RPCRes{feat, ""})
}
