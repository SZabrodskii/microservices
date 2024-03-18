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
	Options  *redis.Options
	TLS      TLSConfig
}

type RedisConfig interface {
	GetDSN() string
	GetPassword() string
	GetDB() int
	GetOptions() *redis.Options
	GetTLSConfig() TLSConfig
}

func (rp *RedisParams) GetDSN() string {
	if rp.DSN == "" {
		rp.DSN = rp.Host + ":" + rp.Port
		//rp.DSN = fmt.Sprintf("host=%s, port=%d", rp.Host, rp.Port)
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
	if rp.Options != nil {
		return rp.Options
	}

	return &redis.Options{}
}

func (rp *RedisParams) GetTLSConfig() TLSConfig {
	return rp.TLS
}

type RedisProvider interface {
	Client() *redis.Client
}
