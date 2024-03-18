package config

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
	if rsp.MasterName == "" {
		return "mymaster"
	}
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
