package providers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	_ "gorm.io/driver/postgres"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type postgresTestCase struct {
	name string
	cfg  *config.PostgreSQLParams
}

func Test_NewPostgresSQLProvider(t *testing.T) {
	logger := zap.NewExample()

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal("cannot connect to Docker - is it running?", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		t.Fatal("Could not connect to Docker", err)
	}

	pool.MaxWait = time.Minute * 2

	matrix := []postgresTestCase{
		{
			name: "without ssl",
			cfg: &config.PostgreSQLParams{
				Username: "admin",
				Password: "Iseestars",
				Database: "testdb",
			},
		},
		{
			name: "with ssl",
			cfg: &config.PostgreSQLParams{
				Username: "admin",
				Password: "Iseestars",
				Database: "testdb",
				TLS: &config.TLSConfig{
					TLSCertificate:     "/etc/postgres/security/server.crt",
					TLSRootCertificate: "/etc/postgres/security/root.crt",
				},
			},
		},
	}

	for _, testCase := range matrix {
		t.Run(testCase.name, func(t *testing.T) {
			resource, err := startPostgresContainer(pool, testCase.cfg.TLS != nil)
			if err != nil {
				t.Fatal("Failed to start PostgreSQL container", err)
			}

			defer func() {
				if err = pool.Purge(resource); err != nil {
					t.Fatal("Could not purge resource", err)
				}
			}()

			hp := strings.Split(resource.GetHostPort("5432/tcp"), ":")
			host := hp[0]
			port, _ := strconv.Atoi(hp[1])
			testCase.cfg.Host, testCase.cfg.Port = host, port
			go func() {
				for {
					tailLogs(pool, resource, context.Background(), os.Stdout, true)
				}
			}()

			assert.Nil(t, err, "Error should be nil")
			assert.Nil(t, waitForPostgres(pool, resource, testCase.cfg.TLS != nil), "Error should be nil")

			provider, err := NewPostgreSQLProvider(testCase.cfg, logger)
			if err != nil {
				t.Fatalf("Error creating PostgreSQL provider: %v", err)
				return
			}

			assert.Nil(t, err, "Error should be nil")
			assert.NotNil(t, provider, "provider should be not nil")
			assert.NotNil(t, provider.db, "DB should not be nil")
		})
	}
}

func startPostgresContainer(pool *dockertest.Pool, useSSL bool) (*dockertest.Resource, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Println("SSL Certificate Paths:", currentPath+"/fixtures/server.crt", currentPath+"/fixtures/root.crt")

	options := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12",
		Env: []string{
			"POSTGRES_USER=admin",
			"POSTGRES_PASSWORD=Iseestars",
			"POSTGRES_DB=testdb",
		},
		Mounts: []string{
			fmt.Sprintf("%s/fixtures/postgresql.conf:/etc/postgresql/config/postgresql.conf", currentPath),
			fmt.Sprintf("%s/fixtures/server.crt:/etc/postgres/security/server.crt", currentPath),
			fmt.Sprintf("%s/fixtures/server.key:/etc/postgres/security/server.key", currentPath),
			fmt.Sprintf("%s/fixtures/root.crt:/etc/postgres/security/root.crt", currentPath),
			fmt.Sprintf("%s/fixtures/pg_hba.conf:/etc/postgresql/config/pg_hba.conf", currentPath),
		},
	}

	return pool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.LogConfig = docker.LogConfig{}
	})
}

func tailLogs(p *dockertest.Pool, c *dockertest.Resource, ctx context.Context, wr io.Writer, follow bool) error {
	opts := docker.LogsOptions{
		Context: ctx,

		Stderr:      true,
		Stdout:      true,
		Follow:      follow,
		Timestamps:  true,
		RawTerminal: true,

		Container: c.Container.ID,

		OutputStream: wr,
	}

	return p.Client.Logs(opts)
}

func waitForPostgres(pool *dockertest.Pool, resource *dockertest.Resource, useSSL bool) error {
	hp := strings.Split(resource.GetHostPort("5432/tcp"), ":")
	host, port := hp[0], hp[1]

	fmt.Printf("Container IP: %s, Port: %s\n", host, port)

	dsn := fmt.Sprintf("host=%s port=%s user=admin password=Iseestars dbname=testdb ", host, port)

	return pool.Retry(func() error {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			return err
		}
		defer db.Close()

		err = db.Ping()
		if err != nil {
			return err
		}
		return nil
	})
}
