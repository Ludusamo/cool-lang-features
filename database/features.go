package database

import (
	"errors"
)

type Feature struct {
	Id          int
	Name        string
	Description string
}

func (d *Database) AddFeature(name string, desc string) (*Feature, error) {
	if _, exists := d.featureMap[name]; exists {
		return nil, errors.New("feature already exists")
	}
	id := len(d.features)
	feat := Feature{id, name, desc}
	d.featureMap[name] = id
	d.features = append(d.features, &feat)

	return &feat, nil
}

func (d *Database) GetFeatures() []*Feature {
	return d.features
}

func (d *Database) GetFeature(id int) (*Feature, error) {
	if id >= len(d.features) {
		return nil, errors.New("feature does not exist")
	}
	feat := d.features[id]
	if feat != nil {
		return feat, nil
	}
	return nil, errors.New("feature does not exist")
}

func (d *Database) DeleteFeature(id int) {
	if id >= len(d.features) {
		return
	}
	feat := d.features[id]
	if feat != nil {
		name := d.features[id].Name
		delete(d.featureMap, name)
	}
}

func (d *Database) ModifyFeature(id int, name string, desc string) (*Feature, error) {
	if id >= len(d.features) {
		return nil, errors.New("feature does not exist")
	}
	feat := d.features[id]
	if feat != nil {
		if feat.Name != name {
			delete(d.featureMap, feat.Name)
			d.featureMap[name] = id
			feat.Name = name
		}
		feat.Description = desc
		return feat, nil
	}
	return nil, errors.New("feature does not exist")
}
