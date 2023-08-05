package integration

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sergicanet9/scv-go-tools/v3/api/utils"
	"github.com/sergicanet9/scv-go-tools/v3/infrastructure"
	"github.com/sergicanet9/scv-go-tools/v3/testutils"
	"github.com/tkudlicka/portflux-api/app/api"
	"github.com/tkudlicka/portflux-api/config"
)

const (
	contentType           = "application/json"
	postgresDBName        = "test-db"
	postgresUser          = "postgres"
	postgresPassword      = "test"
	postgresContainerPort = "5432/tcp"
	postgresDSNEnv        = "postgresDSN"
	jwtSecret             = "eaeBbXUxks"
	nonExpiryToken        = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiYXV0aG9yaXplZCI6dHJ1ZX0.cCKM32os5ROKxeE3IiDWoOyRew9T8puzPUKurPhrDug"
)

// TestMain does the setup before running the tests and the teardown afterwards
func TestMain(m *testing.M) {
	// Uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("could not connect to docker: %s", err)
	}

	postgresResource := setupPostgres(pool)

	// Runs the tests
	code := m.Run()

	if err = pool.Purge(postgresResource); err != nil {
		log.Panicf("could not purge resource: %s", err)
	}
	os.Unsetenv(postgresDSNEnv)

	os.Exit(code)
}

func setupPostgres(pool *dockertest.Pool) *dockertest.Resource {
	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12.6",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", postgresUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", postgresPassword),
			fmt.Sprintf("POSTGRES_DB=%s", postgresDBName),
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Panicf("could not start resource: %s", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", postgresUser, postgresPassword, resource.GetPort(postgresContainerPort), postgresDBName)
	os.Setenv(postgresDSNEnv, dsn)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		_, err = infrastructure.ConnectPostgresDB(context.Background(), dsn)
		return err
	})
	if err != nil {
		log.Panicf("Could not connect to docker: %s", err)
	}

	return resource
}

func Databases(t *testing.T, f func(*testing.T, string), databases ...string) {
	t.Helper()

	// if no databases specified, test is going to run on both
	if len(databases) == 0 {
		databases = []string{"postgres"}
	}

	for _, db := range databases {
		t.Run(db, func(t *testing.T) {
			f(t, db)
		})
	}
}

// New starts a testing instance of the API and returns its config
func New(t *testing.T, database string) config.Config {
	t.Helper()

	cfg, err := testConfig(t, database)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	a := api.New(ctx, cfg)
	run := a.Run(ctx, cancel)
	go run()

	<-time.After(100 * time.Millisecond) // waiting time for letting the API start completely
	return cfg
}

func testConfig(t *testing.T, database string) (c config.Config, err error) {
	c.Version = "Integration tests"
	c.Environment = "Integration tests"
	c.Port = testutils.FreePort(t)
	c.Database = database
	switch database {
	case "postgres":
		c.DSN = os.Getenv(postgresDSNEnv)
	default:
		return config.Config{}, fmt.Errorf("database flag %s not valid", database)
	}

	c.PostgresMigrationsDir = "infrastructure/postgres/migrations"
	c.JWTSecret = jwtSecret
	c.Timeout = utils.Duration{Duration: 30 * time.Second}

	return c, nil
}
