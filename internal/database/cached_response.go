package database

import (
	"context"
)

type CachedResponse struct {
	ID      int    `db:"id" json:"id"`
	URL     string `db:"url" json:"url"`
	Payload []byte `db:"payload" json:"payload"`
}

func (db *DB) InsertCachedResponse(url string, payload []byte) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := `
		INSERT INTO cached_responses (url, payload) VALUES ($1, $2)`

	result, err := db.ExecContext(ctx, query, url, payload)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), err
}

func (db *DB) GetCachedResponse(url string) (*CachedResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	var cachedResponse CachedResponse

	query := `SELECT * FROM cached_responses WHERE url = $1`

	err := db.GetContext(ctx, &cachedResponse, query, url)

	return &cachedResponse, err
}
