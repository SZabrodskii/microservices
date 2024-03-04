package config

import (
	"github.com/go-redis/redis/v8"
)

type RedisParams struct {
	DSN      string
	Host     string
	Port     string
	Password string
	DB       int
	Config   *redis.Options
}

func (rp *RedisParams) GetDSN() string {
	if rp.DSN == "" {
		rp.DSN = rp.Host + ":" + rp.Port
	}
	return rp.DSN
}

func (rp *RedisParams) GetPassword() string {
	return rp.Password
}

func (rp *RedisParams) GetDB() int {
	return rp.DB
}

func (rp *RedisParams) GetOptions() *redis.Options {
	if rp.Config != nil {
		return rp.Config
	}

	return &redis.Options{}
}

type RedisConfig interface {
	GetDSN() string
	GetPassword() string
	GetDB() int
	GetOptions() *redis.Options
}

type RedisProvider interface {
	Client() *redis.Client
}

type RedisSentinelParams struct {
	MasterName       string
	Nodes            []string
	Password         string
	SentinelPassword string
	DB               int
}

func (rsp *RedisSentinelParams) GetMasterName() string {
	return rsp.MasterName
}

func (rsp *RedisSentinelParams) GetNodes() []string {
	return rsp.Nodes
}

func (rsp *RedisSentinelParams) GetPassword() string {
	return rsp.Password
}

func (rsp *RedisSentinelParams) GetSentinelPassword() string {
	return rsp.SentinelPassword
}

func (rsp *RedisSentinelParams) GetDB() int {
	return rsp.DB
}

type RedisSentinelConfig interface {
	GetMasterName() string
	GetNodes() []string
	GetPassword() string
	GetSentinelPassword() string
	GetDB() int
}

type RedisClient struct {
	Options *redis.Options
}

func NewRedisClient(params RedisConfig, options ...*redis.Options) (*RedisClient, error) {
	var opts redis.Options
	if options[0] != nil {
		opts = *options[0]
	}

	opts.Addr = params.GetDSN()

	return &RedisClient{
		Options: &redis.Options{
			Addr:     params.GetDSN(),
			Password: params.GetPassword(),
			DB:       params.GetDB(),
		},
	}, nil
}

func (rc *RedisClient) Client() *redis.Client {
	return redis.NewClient(rc.Options)
}
