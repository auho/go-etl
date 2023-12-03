package dml

import (
	simpledb "github.com/auho/go-simple-db/v2"
)

type manipulation struct {
	st          statement
	prepareFunc func()
}

func (s *manipulation) Sql() string {
	s.prepareFunc()

	return s.st.Query()
}

func (s *manipulation) InsertSql(name string) string {
	s.prepareFunc()

	return s.st.InsertQuery(name)
}

func (s *manipulation) Insert(name string, db *simpledb.SimpleDB) (string, error) {
	_sql := s.InsertSql(name)

	return _sql, db.Exec(_sql).Error
}

func (s *manipulation) InsertWithFieldsSql(name string, fields []string) string {
	s.prepareFunc()

	return s.st.InsertWithFieldsQuery(name, fields)
}

func (s *manipulation) InsertWithField(name string, fields []string, db *simpledb.SimpleDB) (string, error) {
	_sql := s.InsertWithFieldsSql(name, fields)

	return _sql, db.Exec(_sql).Error
}

func (s *manipulation) UpdateSql() string {
	s.prepareFunc()

	return s.st.UpdateQuery()
}

func (s *manipulation) Update(db *simpledb.SimpleDB) (string, error) {
	_sql := s.UpdateSql()

	return _sql, db.Exec(_sql).Error
}

func (s *manipulation) DeleteSql() string {
	s.prepareFunc()

	return s.st.DeleteQuery()
}

func (s *manipulation) Delete(db *simpledb.SimpleDB) (string, error) {
	_sql := s.DeleteSql()

	return _sql, db.Exec(_sql).Error
}
