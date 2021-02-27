package dockertest

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// AdRepository is a repository for ads.
type AdRepository interface {
	Insert(ad Ad) error
	Get(id uuid.UUID) (Ad, error)
	FindAll() []Ad
}

// NewAdRepository returns a concrete repository
func NewAdRepository(db *sql.DB) AdRepository {
	return repo{
		db: db,
	}
}

type repo struct {
	db *sql.DB
}

var _ AdRepository = (*repo)(nil)

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

func (r repo) Get(id uuid.UUID) (Ad, error) {
	sql := `SELECT id, content, start_at, end_at, created
					FROM ad
					WHERE id = $1`

	ad := Ad{}
	row := r.db.QueryRow(sql, id)

	err := row.Scan(
		&ad.ID,
		&ad.Content,
		&ad.StartAt,
		&ad.EndAt,
		&ad.Created,
	)
	return ad, err
}

func (r repo) FindAll() []Ad {
	sql := `SELECT id, content, start_at, end_at, created FROM ad`
	rows, _ := r.db.Query(sql)
	defer rows.Close()

	ads := make([]Ad, 0)
	for rows.Next() {
		ad := Ad{}
		rows.Scan(
			&ad.ID,
			&ad.Content,
			&ad.StartAt,
			&ad.EndAt,
			&ad.Created,
		)
		ads = append(ads, ad)
	}

	return ads
}
