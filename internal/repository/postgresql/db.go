package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(connString string) (*pgxpool.Pool, error) {
	db, dbErr := pgxpool.New(context.Background(), connString)
	if dbErr != nil {
		return nil, dbErr
	}
	pingErr := db.Ping(context.Background())
	if pingErr != nil {
		return nil, pingErr
	}
	fmt.Println("Database is Up and Running!")
	return db, nil
}
