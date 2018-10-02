package server

import (
	"testing"
)

func TestCreateServer(t *testing.T) {
	server := CreateServer()

	if server.db == nil {
		t.Error("database was not initialized")
	}

	if server.router == nil {
		t.Error("router was not initialized")
	}
}
