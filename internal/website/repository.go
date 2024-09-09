package website

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

var ErrNoRows = errors.New("no rows found")

type sqliteWebsite struct {
	uuid      string
	name      string
	url       string
	createdAt string
}

func (w sqliteWebsite) export() (*Website, error) {
	var website Website

	website.Name = w.name
	website.URL = w.url

	uid, err := uuid.Parse(w.uuid)
	if err != nil {
		return nil, err
	}

	website.ID = uid

	d, err := time.Parse(time.RFC3339, w.createdAt)
	if err != nil {
		return nil, err
	}

	website.CreatedAt = d

	return &website, nil
}

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) Create(website Website) (*Website, error) {
	_, err := r.db.Exec(
		"INSERT INTO websites(uuid, name, url, created_at) VALUES(?, ?, ?, ?)",
		website.ID.String(),
		website.Name,
		website.URL,
		website.CreatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return nil, err
	}

	return &website, nil
}

func rowToWebsite(row *sql.Row) (*Website, error) {
	var website sqliteWebsite
	if err := row.Scan(&website.uuid, &website.name, &website.url, &website.createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRows
		}
		return nil, err
	}

	return website.export()
}

func (r *SQLiteRepository) Get(id uuid.UUID) (*Website, error) {
	return rowToWebsite(r.db.QueryRow(
		"SELECT uuid, name, url, created_at FROM websites WHERE uuid = ? LIMIT 1;",
		id.String(),
	))
}

func (r *SQLiteRepository) GetByURL(url string) (*Website, error) {
	return rowToWebsite(r.db.QueryRow(
		"SELECT uuid, name, url, created_at FROM websites WHERE url = ? LIMIT 1;",
		url,
	))
}

func (r *SQLiteRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM websites WHERE uuid = ?;", id.String())
	return err
}

func (r *SQLiteRepository) All() ([]Website, error) {
	rows, err := r.db.Query("SELECT uuid, name, url, created_at FROM websites;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Website
	for rows.Next() {
		var website sqliteWebsite
		if err = rows.Scan(&website.uuid, &website.name, &website.url, &website.createdAt); err != nil {
			return nil, err
		}

		w, err := website.export()
		if err != nil {
			return nil, err
		}

		all = append(all, *w)
	}

	return all, nil
}
