// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package db

import (
	"github.com/hashicorp/go-memdb"
	"github.com/pkg/errors"
)

type DB struct {
	*memdb.MemDB
}

// Create the DB schema
var schema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		"user": &memdb.TableSchema{
			Name: "user",
			Indexes: map[string]*memdb.IndexSchema{
				"id": &memdb.IndexSchema{
					Name:    "id",
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "ID"},
				},
			},
		},
	},
}

func NewDB() (*DB, error) {
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &DB{db}, nil
}
