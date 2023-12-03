package dml

import (
	simpledb "github.com/auho/go-simple-db/v2"
)

type statement struct {
	sr          statementor
	prepareFunc func()
}

func (s *statement) Sql() string {
	s.prepareFunc()

	return s.sr.Query()
}

func (s *statement) InsertSql(name string) string {
	s.prepareFunc()

	return s.sr.InsertQuery(name)
}

func (s *statement) Insert(name string, db *simpledb.SimpleDB) (string, error) {
	_sql := s.InsertSql(name)

	return _sql, db.Exec(_sql).Error
}

func (s *statement) InsertWithFieldsSql(name string, fields []string) string {
	s.prepareFunc()

	return s.sr.InsertWithFieldsQuery(name, fields)
}

func (s *statement) InsertWithField(name string, fields []string, db *simpledb.SimpleDB) (string, error) {
	_sql := s.InsertWithFieldsSql(name, fields)

	return _sql, db.Exec(_sql).Error
}

func (s *statement) UpdateSql() string {
	s.prepareFunc()

	return s.sr.UpdateQuery()
}

func (s *statement) Update(db *simpledb.SimpleDB) (string, error) {
	_sql := s.UpdateSql()

	return _sql, db.Exec(_sql).Error
}

func (s *statement) DeleteSql() string {
	s.prepareFunc()

	return s.sr.DeleteQuery()
}

func (s *statement) Delete(db *simpledb.SimpleDB) (string, error) {
	_sql := s.DeleteSql()

	return _sql, db.Exec(_sql).Error
}
