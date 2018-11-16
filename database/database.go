package database

import (
	"sync"
)

type Database struct {
	featureLocks [1]*sync.RWMutex
	featureMap   map[string]int
	features     []*Feature
}

/** Creates an empty database
 * @return pointer to created database
 */
func CreateDatabase() *Database {
	featureLocks := [1]*sync.RWMutex{&sync.RWMutex{}}
	return &Database{
		featureLocks,
		make(map[string]int),
		make([]*Feature, 0)}
}
