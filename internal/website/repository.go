package website

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"time"
)

var errNoRows = errors.New("no rows found")

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

func (r *SQLiteRepository) Migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS websites (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		url TEXT NOT NULL UNIQUE,
		created_at TEXT
	);
	`

	_, err := r.db.Exec(query)
	return err
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

func (r *SQLiteRepository) Get(id uuid.UUID) (*Website, error) {
	row := r.db.QueryRow("SELECT uuid, name, url, created_at FROM websites WHERE uuid = ?", id.String())

	var website sqliteWebsite
	if err := row.Scan(&website.uuid, &website.name, &website.url, &website.createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errNoRows
		}
		return nil, err
	}

	return website.export()
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
