package database

import (
	"sync"
)

type Database struct {
	featureLocks [10]*sync.RWMutex
	featureMap   map[string]int
	features     []*Feature
}

/** Creates an empty database
 * @return pointer to created database
 */
func CreateDatabase() *Database {
	featureLocks := [10]*sync.RWMutex{}
	for i := 0; i < 10; i++ {
		featureLocks[i] = &sync.RWMutex{}
	}
	return &Database{
		featureLocks,
		make(map[string]int),
		make([]*Feature, 0)}
}
