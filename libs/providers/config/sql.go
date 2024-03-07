package config

type PostgreSQLConfig interface {
	GetDSN() string
	GetUserParams() string
	GetPassword() string
	GetDB() string
}

type PostgreSQLParams struct {
	DSN       string
	Username  string
	Password  string
	Database  string
	SSLConfig *SSLConfig
}

type SSLConfig struct {
	SSLCertificate     string
	SSLRootCertificate string
}

func (p *PostgreSQLParams) GetDSN() string {
	return p.DSN
}

func (p *PostgreSQLParams) GetUserName() string {
	return p.Username
}

func (p *PostgreSQLParams) GetPassword() string {
	return p.Password
}

func (p *PostgreSQLParams) GetDB() string {
	return p.Database
}

func (s *SSLConfig) GetSSLCertificate() string {
	return s.SSLCertificate
}

func (s *SSLConfig) GetSSLRootCertificate() string {
	return s.SSLRootCertificate
}
