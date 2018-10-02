package database

import (
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	db := CreateDatabase()
	if db.features == nil {
		t.Error("database features was not initialized")
	}
}
