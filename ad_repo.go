package dockertest

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AdRepository interface {
	Insert(ad Ad) error
	Delete(id string) error
	Get(id uuid.UUID) (Ad, error)
	FindByStartDate(start time.Time) []Ad
}

func NewAdRepository(db *sql.DB) AdRepository {
	return repo{
		db: db,
	}
}

type repo struct {
	db *sql.DB
}

func (r repo) Insert(ad Ad) error {
	sql := `INSERT INTO ad (id, content, start_at, end_at, created) 
			VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(
		sql,
		ad.ID,
		ad.Content,
		ad.StartAt,
		ad.EndAt,
		time.Now())

	return err
}

func (r repo) Delete(id string) error {
	panic("implement me")
}

func (r repo) Get(id uuid.UUID) (Ad, error) {
	panic("implement me")
}

func (r repo) FindByStartDate(start time.Time) []Ad {
	panic("implement me")
}
