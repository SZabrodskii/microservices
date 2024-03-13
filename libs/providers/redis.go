package providers

import (
	"crypto/tls"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

type RedisClient struct {
	Options *redis.Options
}

func NewRedisClient(params RedisConfig, tlsConfig *tls.Config, options ...*redis.Options) (*RedisClient, error) {
	var opts redis.Options
	if options[0] != nil {
		opts = *options[0]
	}

	opts.Addr = params.GetDSN()
	opts.TLSConfig = tlsConfig

	return &RedisClient{
		Options: &redis.Options{
			Addr:      params.GetDSN(),
			Password:  params.GetPassword(),
			DB:        params.GetDB(),
			TLSConfig: tlsConfig,
		},
	}, nil
}

func (rc *RedisClient) Client() *redis.Client {
	return redis.NewClient(rc.Options)
}

//==================================== End of Experimenting Code =====================================================

type RedisProvider struct {
	client *redis.Client
}

func NewRedisProvider(config config.RedisConfig, tlsConfig *tls.Config) (*RedisProvider, error) {
	opts := config.GetOptions()
	opts.Addr = config.GetDSN()
	opts.TLSConfig = tlsConfig

	if tlsConfig != nil {
		client, err := config.NewRedisClient(config, tlsConfig, opts)
		if err != nil {
			return nil, err
		}

		prvd := &RedisProvider{
			client: client.Client(),
		}

		return prvd, nil
	}

	client := redis.NewClient(opts)

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("error occurred while connecting to the Redis server: %v", err)
	}

	prvd := &RedisProvider{
		client: client,
	}

	return prvd, nil
}

func NewRedisSentinelProvider(config config.RedisSentinelConfig, tlsConfig *tls.Config) (*RedisProvider, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       config.GetMasterName(),
		SentinelAddrs:    config.GetHosts(),
		Password:         config.GetPassword(),
		SentinelPassword: config.GetSentinelPassword(),
		DB:               config.GetDB(),
		TLSConfig:        tlsConfig,
	})

	if tlsConfig != nil {
		client, err := NewRedisClient(config, tlsConfig)
		if err != nil {
			return nil, err
		}

		prvd := &RedisProvider{
			client: client.Client(),
		}

		return prvd, nil
	}

	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("error occurred while connecting to the Redis server: %v", err)
	}

	prvd := &RedisProvider{
		client: client,
	}

	return prvd, nil
}

func (p *RedisProvider) Client() *redis.Client {
	return p.client
}

func ProvideSentinel() fx.Option {
	return fx.Provide(NewRedisSentinelProvider)
}

func ProvideRedis() fx.Option {
	return fx.Provide(NewRedisProvider)
}

func AnnotateSentinel() fx.Option {
	return AnnotateRedis()
}

func AnnotateRedis() fx.Option {
	return fx.Provide(
		fx.Annotate(
			func(client *RedisProvider) *RedisProvider {
				return client
			},
			fx.ResultTags(`name:"redisClient"`),
		),
	)
}

//=========================================== Calling Redis Client ====================================================
//func main() {
//    // Create a new TLS configuration as needed
//    tlsConfig := &tls.Config{
//        // Configure TLS settings as needed
//    }
//
//    // Initialize FX application with the provided modules and options
//    app := fx.New(
//        ProvideSentinel(),
//        ProvideRedis(),
//        AnnotateSentinel(),
//        AnnotateRedis(),
//        // ... other modules
//    )
//
//    // Start the application
//    if err := app.Start(context.Background()); err != nil {
//        log.Fatalf("Error starting application: %v", err)
//    }
//
//    // ...
//
//    // Stop the application
//    if err := app.Stop(context.Background()); err != nil {
//        log.Fatalf("Error stopping application: %v", err)
//    }
//}
