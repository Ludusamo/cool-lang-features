package database

type Database struct {
	featureMap map[string]int
	features   []*Feature
}

/** Creates an empty database
 * @return pointer to created database
 */
func CreateDatabase() *Database {
	return &Database{make(map[string]int), make([]*Feature, 0)}
}
