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
}
