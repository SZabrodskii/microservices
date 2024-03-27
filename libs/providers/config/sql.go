package config

import (
	"fmt"
)

type PostgreSQLConfig interface {
	GetDSN() string
	GetUsername() string
	GetPassword() string
	GetDB() string
	GetTLSConfig() *TLSConfig
}

type PostgreSQLParams struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	TLS      *TLSConfig
}

type TLSConfig struct {
	TLSCertificate     string
	TLSRootCertificate string
}

func (p *PostgreSQLParams) GetDSN() string {
	return fmt.Sprintf("host=%s, port=%d", p.Host, p.Port)
}

func (p *PostgreSQLParams) GetUsername() string {
	return p.Username
}

func (p *PostgreSQLParams) GetPassword() string {
	return p.Password
}

func (p *PostgreSQLParams) GetDB() string {
	return p.Database
}

func (p *PostgreSQLParams) GetTLSConfig() *TLSConfig {
	return p.TLS
}

func (s *TLSConfig) GetTLSCertificate() string {
	return s.TLSCertificate
}

func (s *TLSConfig) GetTLSRootCertificate() string {
	return s.TLSRootCertificate
}
