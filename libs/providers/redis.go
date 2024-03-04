package providers

import (
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
)

type RedisProvider struct {
	client *redis.Client
}

func NewRedisProvider(config config.RedisConfig) (*RedisProvider, error) {
	opts := config.GetOptions()
	opts.Addr = config.GetDSN()

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

func NewRedisSentinelProvider(config config.RedisSentinelConfig) (*RedisProvider, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       config.GetMasterName(),
		SentinelAddrs:    config.GetNodes(),
		Password:         config.GetPassword(),
		SentinelPassword: config.GetSentinelPassword(),
		DB:               config.GetDB(),
	})

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

func AnnotateSentinel() fx.Option {
	return AnnotateRedis()
}
