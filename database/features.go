package database

import (
	"errors"
)

type Feature struct {
	id          int
	name        string
	description string
}

func (d *Database) AddFeature(name string, desc string) (*Feature, error) {
	if _, exists := d.features[name]; exists {
		return nil, errors.New("feature already exists")
	}
	feat := Feature{len(d.features), name, desc}
	d.features[name] = &feat
	return &feat, nil
}

func (d *Database) DeleteFeature(name string) {
    delete(d.features, name)
}
