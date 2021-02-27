package dockertest_test

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ayubmalik/dockertest"
)

var db *sql.DB

func OpenDB(host string, port, user, pwd, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func TestMain(m *testing.M) {
	_db, err := OpenDB("localhost", "5432", "postgres", "", "dockertest")
	if err != nil {
		panic(err)
	}
	db = _db
	defer db.Close()
	os.Exit(m.Run())
}

func InitDB(t *testing.T) {
	f, err := os.Open("migrations/001-create-db.sql")
	must(t, err)
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	must(t, err)

	_, err = db.Exec(string(buf))
	must(t, err)
}

func TestAdRepoInsert(t *testing.T) {
	InitDB(t)
	repo := dockertest.NewAdRepository(db)

	now := time.Now()
	start := now.AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 1)

	ad := dockertest.NewAd("my ad content", start, end)

	err := repo.Insert(ad)
	must(t, err)

	var (
		id             uuid.UUID
		content        string
		startAt, endAt time.Time
	)

	row := db.QueryRow("select id, content, start_at, end_at from ad where id = $1", ad.ID)
	err = row.Scan(&id, &content, &startAt, &endAt)
	must(t, err)
	assert(t, id, ad.ID)
}

func assert(t *testing.T, got, want interface{}) {
	if got != want {
		t.Errorf("got %v wanted %v", got, want)
	}
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
