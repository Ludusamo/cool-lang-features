package database

import (
	"errors"
)

type Feature struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

/** Calculates the bucket that an id is in
 * @lhs database pointer
 * @param id the id to find the bucket for
 * @return the bucket number
 */
func (d *Database) calcBucket(id int) int {
	bucketSize := len(d.features) / len(d.featureLocks)
	return id / bucketSize
}

/** Adds an entry for a feature in the database
 * @param name string identifier for feature
 * @param desc text description of the feature
 * @return pointer to added feature, an error if it could not be added
 */
func (d *Database) AddFeature(name string, desc string) (*Feature, error) {
	for _, lock := range d.featureLocks {
		lock.Lock()
		defer lock.Unlock()
	}
	if _, exists := d.featureMap[name]; exists {
		return nil, errors.New("feature already exists")
	}
	id := len(d.features)
	feat := Feature{id, name, desc}
	d.featureMap[name] = id
	d.features = append(d.features, &feat)

	return &feat, nil
}

/** Retrieves an array of all features in the database
 * @return array of all features
 */
func (d *Database) GetFeatures() []*Feature {
	for _, lock := range d.featureLocks {
		lock.RLock()
		defer lock.RUnlock()
	}
	cpy := make([]*Feature, len(d.features))
	copy(cpy, d.features)
	return cpy
}

/** Retrieve a specific feature from its identifier
 * @param id integer identifier of feature
 * @return feature with the given identifier, error if it doesn't exist
 */
func (d *Database) GetFeature(id int) (*Feature, error) {
	if id >= len(d.features) {
		return nil, errors.New("feature does not exist")
	}
	lock := d.featureLocks[d.calcBucket(id)]
	lock.RLock()
	defer lock.RUnlock()
	feat := d.features[id]
	if feat != nil {
		return feat, nil
	}
	return nil, errors.New("feature does not exist")
}

/** Delete a feature given its identifier
 * @param id integer identifier of feature
 */
func (d *Database) DeleteFeature(id int) {
	if id >= len(d.features) {
		return
	}
	lock := d.featureLocks[d.calcBucket(id)]
	lock.Lock()
	defer lock.Unlock()
	feat := d.features[id]
	if feat != nil {
		name := d.features[id].Name
		delete(d.featureMap, name)
		d.features[id] = nil
	}
}

/** Modify a feature given its identifier, the id itself cannot be changed
 * @param id integer identifier of feature
 * @param name new string identification
 * @param desc new text description of feature
 * @return modified feature, error if it doesn't exist
 */
func (d *Database) ModifyFeature(id int, name string, desc string) (*Feature, error) {
	if id >= len(d.features) {
		return nil, errors.New("feature does not exist")
	}
	lock := d.featureLocks[d.calcBucket(id)]
	lock.Lock()
	defer lock.Unlock()
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
