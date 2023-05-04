package repository

import (
	"context"
	"zentotem/pkg/postgres"

	_ "github.com/lib/pq"
)

type UserStorage struct {
	Db *postgres.PG
}

func NewUserStorage(db *postgres.PG) *UserStorage {
	return &UserStorage{Db: db}
}

func (s *UserStorage) CreateUser(ctx context.Context, name string, age uint32) (userId uint64, err error) {
	tx := s.Db.ExtractTx(ctx)
	var id uint64
	query := `
	INSERT INTO user 
	(name, age) 
	VALUES ($1,$2) returning id`
	err = tx.GetContext(ctx, &id, query, name, age)
	if err != nil {
		return 0, err
	}
	return id, nil
}
