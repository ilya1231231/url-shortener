package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"url-shortender/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	//Здесь храниится имя функции, которое покажется при ошибке
	const fn = "storage.sqlite.New"
	handle := func(fn string, err error) (*Storage, error) {
		//врапим ошибку, чтобы было понятно где она возникла
		//Errorf позволяет нам использовать функции форматирования для создания описательных сообщений об ошибках.

		//!!! Важно
		//С глаголом аннотации %w возвращаемая ошибка анврапится в обычную ошибку
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return handle(fn, err)
	}
	//stmt - statement. sql, который можно выполнить позже.
	//Обычно в тут должен храниться запрос с параметрами, которые можно подставить позже и выполнить запрос
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url(
    			id INTEGER PRIMARY KEY, 
				alias TEXT NOT NULL UNIQUE, 
				url TEXT NOT NULL);
				CREATE INDEX IF NOT EXISTS idx_alias ON alias`)

	if err != nil {
		return handle(fn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return handle(fn, err)
	}

	return &Storage{db: db}, nil
}

// изучить подробнее
func (s *Storage) CreateUrl(url string, alias string) (int64, error) {
	const fn = "storage.sqlite.CreateUrl"
	handle := func(fn string, err error) (int64, error) {
		return 0, fmt.Errorf("%s:%w", fn, err)
	}

	stmt, err := s.db.Prepare(`INSERT INTO url (url, alias) VALUES (?, ?)`)
	if err != nil {
		return handle(fn, err)
	}
	res, err := stmt.Exec(url, alias)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return handle(fn, storage.ErrURLExists)
		}
		return handle(fn, err)
	}

	lastInsertId, err := res.LastInsertId()

	if err != nil {
		return handle(fn, err)
	}
	return lastInsertId, nil
}

type Url struct {
	ID    int64
	Alias string
	Url   string
}

func (s *Storage) GetUrl(alias string) (*Url, error) {
	var UrlData Url
	const fn = "storage.sqlite.GetUrl"
	handle := func(fn string, err error) (*Url, error) {
		return nil, fmt.Errorf("%s:%w", fn, err)
	}

	stmt, err := s.db.Prepare(`SELECT * FROM url WHERE alias = ?`)
	if err != nil {
		return handle(fn, err)
	}

	err = stmt.QueryRow(alias).Scan(&UrlData.ID, &UrlData.Alias, &UrlData.Url)
	if err != nil {
		return handle(fn, err)
	}
	return &UrlData, nil
}
