package database

import (
	"testing"
)

func TestAddFeature(t *testing.T) {
	db := Database{make(map[string]*Feature)}
	name := "Test Feature"
	description := "This is an example description."

	feat, err := db.AddFeature(name, description)
	if err != nil {
		t.Fatal(err)
	}

	newFeature, exists := db.features[name]
	if !exists {
		t.Fatal("feature doesn't exist in the database")
	}

	if feat != newFeature {
		t.Error("returned feature does not equal retrieved feature")
	}

	if newFeature.id != 0 {
		t.Errorf("incorrect id set: received %v, expected 0", newFeature.id)
	}

	if newFeature.name != name {
		t.Errorf("incorrect name set: received %v, expected %v",
			newFeature.name,
			name)
	}

	if newFeature.description != description {
		t.Errorf("incorrect description set: received %v, expected %v",
			newFeature.description,
			description)
	}

	_, dupErr := db.AddFeature(name, description)
	if dupErr == nil || dupErr.Error() != "feature already exists" {
		t.Error("adding a duplicate feature does not fail")
	}

	other, _ := db.AddFeature("other", description)
	if other.id != 1 {
		t.Errorf("incorrect id set for new entry: received %v, expected 1",
			newFeature.id)
	}
}

func TestGetFeature(t *testing.T) {
    db := Database{make(map[string]*Feature)}
    name := "Test"
    db.features[name] = &Feature{0, name, ""}

    feat, err := db.GetFeature(name)
    if err != nil {
        t.Fatal(err)
    }

    if feat.name != name || feat.description != "" {
        t.Error("feature retrieved does not match expected output")
    }

    _, dneErr := db.GetFeature("DNE")
    if dneErr == nil || dneErr.Error() != "feature does not exist" {
        t.Error("retrieving a feature that does not exist succeeded")
    }
}

func TestRemoveFeature(t *testing.T) {
	db := Database{make(map[string]*Feature)}

	// Create mock feature
	name := "test"
	db.features[name] = &Feature{0, name, ""}

	db.DeleteFeature(name)
	_, exists := db.features[name]
	if exists {
		t.Error("feature still exists after deleting")
	}
}

func TestModifyFeature(t *testing.T) {
	db := Database{make(map[string]*Feature)}
	name := "Test"
	db.features[name] = &Feature{0, name, ""}

	// Modify description
	newDesc := "New Description"
	modifiedFeat, err := db.ModifyFeature(name, newDesc)
	if err != nil {
		t.Fatal(err)
	}
	if modifiedFeat != db.features[name] {
		t.Error("returned features does not match feature in database")
	}
	if modifiedFeat.description != newDesc {
		t.Errorf("incorrect description set: received %v, expected %v",
			modifiedFeat.description,
			newDesc)
	}

	// Trying to modify a feature that doesn't exist
	_, dneErr := db.ModifyFeature("DNE", "Some Description")
	if dneErr == nil || dneErr.Error() != "feature does not exist" {
		t.Error("error did not come back for a feature that does not exist")
	}
}
