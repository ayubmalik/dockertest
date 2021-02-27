package dockertestspike_test

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	dts "github.com/ayubmalik/dockertestspike"
)

// use docketest by passing -dt flag e.g. go test -dt
var dt = flag.Bool("dt", false, "Use dockertest container rather than external postgres")

var pool *dockertest.Pool

func TestMain(m *testing.M) {
	var err error
	flag.Parse()
	if *dt {
		fmt.Println("Using docker test container")
		pool, err = dockertest.NewPool("")
		if err != nil {
			panic(err)
		}
	}

	os.Exit(m.Run())
}

func TestAdRepoInsert(t *testing.T) {
	t.Parallel()
	db := openDB(t)
	repo := dts.NewAdRepository(db)

	now := time.Now()
	start := now.AddDate(0, 0, 1)
	end := start.AddDate(0, 0, 1)

	ad := dts.NewAd("my ad content", start, end)

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
	t.Parallel()
	db := openDB(t)
	repo := dts.NewAdRepository(db)

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
	t.Parallel()
	db := openDB(t)
	repo := dts.NewAdRepository(db)

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

func openDB(t *testing.T) *sql.DB {
	port := "5432"

	if *dt {
		resource, err := pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "13.2-alpine",
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=password",
				"POSTGRES_DB=dockertest",
			},
			ExposedPorts: []string{port},
		}, func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})

		if err != nil {
			t.Fatalf("Could not start resource: %s", err)
		}

		t.Cleanup(func() {
			err := pool.Purge(resource)
			if err != nil {
				t.Logf("Could not purge resource: %s", err)
			}
		})

		port = resource.GetPort(fmt.Sprintf("%s/tcp", port))
		if err := pool.Retry(func() error {
			_db, _err := sql.Open("postgres", dsn(port))
			if _err != nil {
				return _err
			}
			defer _db.Close()
			return _db.Ping()
		}); err != nil {
			log.Fatalf("Could not connect to docker/postgres after retry %s", err)
		}

		fmt.Println("postgres ready on port:", port)
	}

	db, err := sql.Open("postgres", dsn(port))
	if err != nil {
		t.Fatalf("Could not connect DB: %s", err)
	}

	dbUp(t, db)

	return db
}

func dsn(port string) string {
	return fmt.Sprintf("host=localhost port=%s user=postgres password=password dbname=dockertest sslmode=disable", port)
}

func dbUp(t *testing.T, db *sql.DB) {
	f, err := os.Open("migrations/001-create-db.sql")
	must(t, err)
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	must(t, err)

	_, err = db.Exec(string(buf))
	must(t, err)
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assert(t *testing.T, got, want interface{}) {
	if got != want {
		t.Errorf("got %v wanted %v", got, want)
	}
}
