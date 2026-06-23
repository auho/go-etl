package dml

import (
	simpledb "github.com/auho/go-simple-db/v2"
)

type manipulation struct {
	st          statement
	prepareFunc func()
}

func (s *manipulation) SQL() string {
	s.prepareFunc()

	return s.st.Query()
}

func (s *manipulation) InsertSQL(name string) string {
	s.prepareFunc()

	return s.st.InsertQuery(name)
}

func (s *manipulation) Insert(name string, db *simpledb.SimpleDB) (string, error) {
	_sql := s.InsertSQL(name)

	return _sql, db.GormDB().Exec(_sql).Error
}

func (s *manipulation) InsertWithFieldsSQL(name string, fields []string) string {
	s.prepareFunc()

	return s.st.InsertWithFieldsQuery(name, fields)
}

func (s *manipulation) InsertWithFields(name string, fields []string, db *simpledb.SimpleDB) (string, error) {
	_sql := s.InsertWithFieldsSQL(name, fields)

	return _sql, db.GormDB().Exec(_sql).Error
}

func (s *manipulation) UpdateSQL() string {
	s.prepareFunc()

	return s.st.UpdateQuery()
}

func (s *manipulation) Update(db *simpledb.SimpleDB) (string, error) {
	_sql := s.UpdateSQL()

	return _sql, db.GormDB().Exec(_sql).Error
}

func (s *manipulation) DeleteSQL() string {
	s.prepareFunc()

	return s.st.DeleteQuery()
}

func (s *manipulation) Delete(db *simpledb.SimpleDB) (string, error) {
	_sql := s.DeleteSQL()

	return _sql, db.GormDB().Exec(_sql).Error
}
