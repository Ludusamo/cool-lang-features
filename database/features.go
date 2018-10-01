package database

import (
    "errors"
)

type Feature struct {
    id int
    name string
    description string
}

func (d *Database) AddFeature(name string, desc string) (*Feature, error) {
    return nil, errors.New("feature could not be created")
}
