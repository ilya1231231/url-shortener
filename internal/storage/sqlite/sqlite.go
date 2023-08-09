package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	//Здесь храниится имя функции, которое покажется при ошибке
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		//врапим ошибку, чтобы было понятно где она возникла
		//Errorf позволяет нам использовать функции форматирования для создания описательных сообщений об ошибках.

		//!!! Важно
		//С глаголом аннотации %w возвращаемая ошибка анврапится в обычную ошибку
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//изучить лучше
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url(
    			id INTEGER PRIMARY KEY, 
				alias TEXT NOT NULL UNIQUE, 
				url TEXT NOT NULL);
				CREATE INDEX IF NOT EXISTS idx_alias ON alias`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
