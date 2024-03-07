package providers

import (
	"database/sql"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"log"
	"testing"

	_ "gorm.io/driver/postgres"
)

func Test_NewPostgresSQLProvider(t *testing.T) {
	logger := zap.NewExample()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("cannot connect to Docker - is it running? %s", err)
	}
	resource, err := startPostgresContainer(pool)
	defer resource.Close()

	assert.Nil(t, err, "Error should be nil")
	assert.Nil(t, pool.Purge(resource), "Error should be nil")
	assert.Nil(t, waitForPostgres(pool, resource), "Error should be nil")

	tempConfig := config.PostgreSQLParams{
		Username: "admin",
		Password: "Iseestars",
		Database: "testdb",
	}

	provider, err := NewPostgreSQLProvider(tempConfig, logger)
	assert.Nil(t, err, "Error should ne nil")
	assert.NotNil(t, provider, "provider should be not nil")
	assert.NotNil(t, provider.db, "DB should not be nil")

}

func startPostgresContainer(pool *dockertest.Pool) (*dockertest.Resource, error) {
	options := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12-alpine",
		Env: []string{
			"POSTGRES_USER=admin",
			"POSTGRES_PASSWORD=Iseestars",
			"POSTGRES_DB=testdb",
		},
		Cmd: []string{"postgres", "-c", "config_file=/etc/postgresql/config/postgresql.conf", "-c", "hba_file=/etc/postgresql/config/pg_hba.conf"},
		Mounts: []string{
			"fixtures/postgresql.conf:/etc/postgresql/config/postgresql.conf",
			"fixtures/server.crt:/etc/postgres/security/server.crt",
			"fixtures/server.key:/etc/postgres/security/server.key",
			"fixtures/root.crt:/etc/postgres/security/root.crt",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "0.0.0.0"}},
		},
	}

	return pool.RunWithOptions(options)
}

func waitForPostgres(pool *dockertest.Pool, resource *dockertest.Resource) error {
	return pool.Retry(func() error {
		db, err := sql.Open("pgx", fmt.Sprintf("host=localhost port=%s user=admin password=Iseestars dbname=testdb sslmode=disable", resource.GetPort("5432/tcp")))
		if err != nil {
			return err
		}
		defer db.Close()

		return db.Ping()
	})
}
