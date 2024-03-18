package providers

import (
	_ "crypto/tls"
	"database/sql"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"github.com/alicebob/miniredis/v2"
	_ "github.com/alicebob/miniredis/v2"
	_ "github.com/go-redis/redis/v8"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"os"
	_ "strconv"
	"strings"
	"testing"
)

func TestNewRedisProviderSSLOff(t *testing.T) {
	miniRedis := miniredis.NewMiniRedis()
	if err := miniRedis.Start(); err != nil {
		t.Fatal("failed to start mini Redis server:", err)
	}
	defer miniRedis.Close()

	cfg := &config.RedisParams{
		DSN: miniRedis.Addr(),
		DB:  0,
	}

	t.Run("should connect without tls", func(t *testing.T) {
		_, err := NewRedisProvider(cfg)
		if err != nil {
			t.Errorf("Failed to create Redis provider: %v", err)
		}

	})

}

func TestNewRedisProviderSSLOn(t *testing.T) {
	pool, err := dockertest.NewTLSPool("", "/etc/libs/providers/fixtures/redis.crt")
	if err != nil {
		t.Fatalf("%v: cannot connect to Docker - is it running?", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %v", err)
	}

	resource, err := startRedisContainer(pool)
	if err != nil {
		t.Fatal("Failed to start redisTLS container", err)
	}
	defer func() {
		if err = pool.Purge(resource); err != nil {
			t.Fatal("Could not purge resource", err)
		}
	}()

	ip := resource.GetBoundIP("6379/tcp")
	port := resource.GetPort("6379/tcp")

	testCfg := &config.RedisParams{
		DSN: fmt.Sprintf("redis://%s:%s", ip, port),
		DB:  0,
		TLS: &config.TLSParams{
			TLSKey:             "/etc/libs/providers/fixtures/redis.key",
			TLSCertificate:     "/etc/libs/providers/fixtures/redis.crt",
			TLSRootCertificate: "/etc/libs/providers/fixtures/redis.csr",
		},
	}

	assert.Nil(t, waitForRedis(pool, resource), "Error should be nil")

	t.Run("should connect with TLS", func(t *testing.T) {
		_, err := NewRedisProvider(testCfg)
		if err != nil {
			t.Errorf("cannot connect to DB with TLS mode: %v", err)
		}

	})

	//hp := strings.Split(resource.GetHostPort("5432/tcp"), ":")
	//host := hp[0]
	//port, _ := strconv.Atoi(hp[1])
	//test.cfg.Host, test.cfg.Port = host, port
	//go func() {
	//	for {
	//		tailLogs(pool, resource, context.Background(), os.Stdout, true)
	//	}
	//}()

}

func startRedisContainer(pool *dockertest.Pool) (*dockertest.Resource, error) {

	currentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Println("SSL Certificate Paths:", currentPath+"/fixtures/redis.crt", currentPath+"/fixtures/redis.key")

	options := &dockertest.RunOptions{
		Repository: "redis",
		Tag:        "12",
		Env: []string{
			"REDIS_USER=admin",
			"REDIS_PASSWORD=Iseestars",
			"REDIS_DB=testdb",
		},
		Mounts: []string{
			fmt.Sprintf("%s/fixtures/redis.crt:/etc/redis/security/redis.crt", currentPath),
			fmt.Sprintf("%s/fixtures/redis.csr:/etc/redis/security/redis.csr", currentPath),
			fmt.Sprintf("%s/fixtures/redis.key:/etc/redis/security/redis.key", currentPath),
			fmt.Sprintf("%s/fixtures/redis.pem:/etc/redis/security/redis.pem", currentPath),
		},
	}

	return pool.RunWithOptions(options, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		config.LogConfig = docker.LogConfig{}
	})

}

func waitForRedis(pool *dockertest.Pool, resource *dockertest.Resource) error {
	hp := strings.Split(resource.GetHostPort("6379/tcp"), ":")
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
