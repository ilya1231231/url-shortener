package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	//Здесь храниится имя функции, которое покажется при ошибке
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		//врапим ошибку, чтобы было понятно где она возникла
		//Errorf позволяет нам использовать функции форматирования для создания описательных сообщений об ошибках.

		//!!! Важно
		//С глаголом аннотации %w возвращаемая ошибка анврапится в обычную ошибку
		return nil, fmt.Errorf("%s: %w", op, err)
	}
}
