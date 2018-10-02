package database

type Database struct {
	features map[string]*Feature
}

func CreateDatabase() *Database {
	return &Database{make(map[string]*Feature)}
}
