package dockertest_test

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ayubmalik/dockertest"
)

// use docketest by passing -dt flag e.g. go test -dt
var dt = flag.Bool("dt", false, "Use dockertest container rather than external postgres")

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
	flag.Parse()
	if *dt {
		fmt.Println("Using docker test container")
	} else {
		fmt.Println("Using external postgres")
	}

	_db, err := OpenDB("localhost", "5432", "postgres", "password", "dockertest")
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

func TestAdRepoGet(t *testing.T) {
	InitDB(t)
	repo := dockertest.NewAdRepository(db)

	id := uuid.New()
	now := time.Now()

	_, err := db.Exec(
		"insert into ad(id, content, start_at, end_at, created) values($1, $2, $3, $4, $5)",
		id,
		"hello",
		now,
		now,
		now,
	)
	must(t, err)

	ad, err := repo.Get(id)
	must(t, err)
	assert(t, ad.ID, id)
	assert(t, ad.Content, "hello")
	assert(t, ad.StartAt.Format(time.RFC3339), now.Format(time.RFC3339))
	assert(t, ad.EndAt.Format(time.RFC3339), now.Format(time.RFC3339))
	assert(t, ad.Created.Format(time.RFC3339), now.Format(time.RFC3339))

}

func TestAdRepoFindAll(t *testing.T) {
	InitDB(t)
	repo := dockertest.NewAdRepository(db)

	for i := 0; i < 100; i++ {
		id := uuid.New()
		now := time.Now()
		_, err := db.Exec(
			"insert into ad(id, content, start_at, end_at, created) values($1, $2, $3, $4, $5)",
			id,
			"hello",
			now,
			now,
			now,
		)
		must(t, err)
	}

	ads := repo.FindAll()
	assert(t, len(ads), 100)
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
