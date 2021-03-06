package database

import (
	"testing"
)

func TestAddFeature(t *testing.T) {
	db := Database{make(map[string]int), make([]*Feature, 0)}
	name := "Test Feature"
	description := "This is an example description."

	feat, err := db.AddFeature(name, description)
	if err != nil {
		t.Fatal(err)
	}

	newFeature := db.features[0]
	if newFeature == nil {
		t.Fatal("feature doesn't exist in the database")
	}

	if feat != newFeature {
		t.Error("returned feature does not equal retrieved feature")
	}

	if newFeature.Id != 0 {
		t.Errorf("incorrect id set: received %v, expected 0", newFeature.Id)
	}

	if newFeature.Name != name {
		t.Errorf("incorrect name set: received %v, expected %v",
			newFeature.Name,
			name)
	}

	if newFeature.Description != description {
		t.Errorf("incorrect description set: received %v, expected %v",
			newFeature.Description,
			description)
	}

	_, dupErr := db.AddFeature(name, description)
	if dupErr == nil || dupErr.Error() != "feature already exists" {
		t.Error("adding a duplicate feature does not fail")
	}

	other, _ := db.AddFeature("other", description)
	if other.Id != 1 {
		t.Errorf("incorrect id set for new entry: received %v, expected 1",
			newFeature.Id)
	}
}

func TestGetFeatures(t *testing.T) {
	db := Database{make(map[string]int), make([]*Feature, 2)}
	db.features[0] = &Feature{0, "First", ""}
	db.features[1] = &Feature{1, "Second", ""}

	features := db.GetFeatures()
	if features == nil {
		t.Fatal("no features retrieved")
	}
	if len(features) != 2 {
		t.Fatal("features does not have the right amount of elements")
	}
	if features[0].Name != "First" {
		t.Errorf("expected first element name to be %v, received %v",
			"First",
			features[0].Name)
	}
	if features[1].Name != "Second" {
		t.Errorf("expected first element name to be %v, received %v",
			"Second",
			features[1].Name)
	}
}

func TestGetFeature(t *testing.T) {
	db := Database{make(map[string]int), make([]*Feature, 1)}
	name := "Test"
	db.features[0] = &Feature{0, name, ""}

	feat, err := db.GetFeature(0)
	if err != nil {
		t.Fatal(err)
	}

	if feat.Name != name || feat.Description != "" {
		t.Error("feature retrieved does not match expected output")
	}

	_, dneErr := db.GetFeature(1)
	if dneErr == nil || dneErr.Error() != "feature does not exist" {
		t.Error("retrieving a feature that does not exist succeeded")
	}
}

func TestDeleteFeature(t *testing.T) {
	db := Database{make(map[string]int), make([]*Feature, 1)}

	// Create mock feature
	name := "test"
	db.featureMap[name] = 0
	db.features[0] = &Feature{0, name, ""}

	db.DeleteFeature(0)
	feat := db.features[0]
	_, exists := db.featureMap[name]
	if feat != nil || exists {
		t.Error("feature still exists after deleting")
	}

	db.DeleteFeature(100)
}

func TestModifyFeature(t *testing.T) {
	db := Database{make(map[string]int), make([]*Feature, 1)}
	name := "Test"
	db.featureMap[name] = 0
	db.features[0] = &Feature{0, name, ""}

	// Modify description
	newDesc := "New Description"
	modifiedFeat, err := db.ModifyFeature(0, name, newDesc)
	if err != nil {
		t.Fatal(err)
	}
	if modifiedFeat != db.features[0] {
		t.Error("returned features does not match feature in database")
	}
	if modifiedFeat.Description != newDesc {
		t.Errorf("incorrect description set: received %v, expected %v",
			modifiedFeat.Description,
			newDesc)
	}

	// Modify name
	newName := "New Name"
	modifiedFeat, err = db.ModifyFeature(0, newName, newDesc)
	if err != nil {
		t.Fatal(err)
	}
	if modifiedFeat != db.features[0] {
		t.Error("returned features does not match feature in database")
	}
	if modifiedFeat.Name != newName {
		t.Errorf("incorrect name set: received %v, expected %v",
			modifiedFeat.Name,
			newName)
	}

	// Trying to modify a feature that doesn't exist
	_, dneErr := db.ModifyFeature(3, "DNE", "Some Description")
	if dneErr == nil || dneErr.Error() != "feature does not exist" {
		t.Error("error did not come back for a feature that does not exist")
	}
}
