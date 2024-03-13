package config

import (
	_ "github.com/SZabrodskii/microservices/libs/providers"
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

// ======================================== Redis Sentinel ============================================================

type RedisSentinelParams struct {
	MasterName       string
	Hosts            []string
	RedisPassword    string
	SentinelPassword string
	DB               int
}

type RedisSentinelConfig interface {
	GetMasterName() string
	GetHosts() []string
	GetPassword() string
	GetSentinelPassword() string
	GetDB() int
}

func (rsp *RedisSentinelParams) GetMasterName() string {
	return rsp.MasterName
}

func (rsp *RedisSentinelParams) GetHosts() []string {
	return rsp.Hosts
}

func (rsp *RedisSentinelParams) GetPassword() string {
	return rsp.RedisPassword
}

func (rsp *RedisSentinelParams) GetSentinelPassword() string {
	if rsp.SentinelPassword == "" {
		return rsp.RedisPassword
	}
	return rsp.SentinelPassword
}

func (rsp *RedisSentinelParams) GetDB() int {
	return rsp.DB
}

//================================================ Redis Client ======================================================

//================================================= TLS Config =======================================================

type TLSConfig struct {
}
