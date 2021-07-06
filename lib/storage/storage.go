package storage

type Source interface {
}

type Target interface {
}

type DbSource struct {
	sqlTemplate    string
	maxSqlTemplate string
}

type DbTarget struct {
	maxConcurrent  int
	size           int
	driver         string
	dsn            string
	scheme         string
	table          string
	sqlTemplate    string
	maxSqlTemplate string
}
