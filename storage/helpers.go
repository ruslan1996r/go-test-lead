package storage

import "os"

type SQLHelpersReader interface {
	ReadSQLFile(filepath string) (string, error)
}

type SQLHelpers struct{}

func NewSQLHelper() *SQLHelpers {
	return &SQLHelpers{}
}

func (s *SQLHelpers) ReadSQLFile(filepath string) (string, error) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
