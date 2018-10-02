package database

type Database struct {
	featureMap map[string]int
	features   []*Feature
}

func CreateDatabase() *Database {
	return &Database{make(map[string]int), make([]*Feature, 0)}
}
