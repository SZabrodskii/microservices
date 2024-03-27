package providers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	_ "github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

type PostgreSQLProvider struct {
	db *gorm.DB
}

type PostgreSQLProviderOptions struct {
	DB      *gorm.DB
	Options *gorm.Config
}

type zapLogger struct {
	log *zap.Logger
}

func (z *zapLogger) Info(_ context.Context, msg string, fields ...interface{}) {
	z.log.Info(msg, zap.Any("fields", fields))
}

func (z *zapLogger) Warn(_ context.Context, msg string, fields ...interface{}) {
	z.log.Warn(msg, zap.Any("fields", fields))
}

func (z *zapLogger) Error(_ context.Context, msg string, fields ...interface{}) {
	z.log.Error(msg, zap.Any("fields", fields))
}

func (z *zapLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	passedTime := time.Since(begin)
	sql, rowsAffected := fc()

	z.log.Debug("SQL trace",
		zap.String("sql", sql),
		zap.Int64("rows_affected", rowsAffected),
		zap.Error(err),
		zap.Duration("elapsed", passedTime),
	)
}

func (z *zapLogger) LogMode(level logger.LogLevel) logger.Interface {

	return z
}

func newZapLogger(log *zap.Logger) *zapLogger {
	return &zapLogger{
		log: log,
	}
}

func NewPostgreSQLProvider(cfg config.PostgreSQLConfig, log *zap.Logger, options ...*PostgreSQLProviderOptions) (*PostgreSQLProvider, error) {
	var (
		dialector gorm.Dialector
		err       error
		conn      *gorm.DB
		opts      *gorm.Config
	)

	if len(options) > 0 {
		if options[0].DB != nil {
			conn = options[0].DB
		}
		if options[0].Options != nil {
			opts = options[0].Options
		}
	}

	dsn := fmt.Sprintf("%s user=%s password=%s dbname=%s", cfg.GetDSN(), cfg.GetUsername(), cfg.GetPassword(), cfg.GetDB())

	dialector = postgres.Open(dsn)

	if opts == nil {
		opts = &gorm.Config{}
	}

	opts.Logger = newZapLogger(log)
	conn, err = gorm.Open(dialector, opts)
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to connect the db: %v", err)
	}

	sql, err := conn.DB()
	if err != nil {
		return nil, fmt.Errorf("an error occured while trying to getting the db: %v", err)
	}

	if err = sql.Ping(); err != nil {
		return nil, fmt.Errorf("canot ping db: %v", err)
	}

	return &PostgreSQLProvider{
		db: conn,
	}, nil
}

func (p *PostgreSQLProvider) getSSLCertificate(certificatePath, rootCertificatePath string) (*tls.Config, error) {
	certificate, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading SSL certificate: %v", err)
	}

	rootCertificate, err := os.ReadFile(rootCertificatePath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading SSL root certificate: %v", err)
	}

	parsedCertificate, err := x509.ParseCertificate(certificate)
	if err != nil {
		return nil, fmt.Errorf("error occured while parsing certificate: %v", err)
	}

	rootCertificates, err := x509.ParseCertificates(rootCertificate)
	if err != nil {

		return nil, fmt.Errorf("Error parsing root certificate: %v", err)
	}

	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}

	for _, rootCertificate := range rootCertificates {
		tlsConfig.RootCAs.AddCert(rootCertificate)
	}

	opts := x509.VerifyOptions{
		Roots:         tlsConfig.RootCAs,
		Intermediates: x509.NewCertPool()}

	_, err = parsedCertificate.Verify(opts)
	if err != nil {
		return nil, fmt.Errorf("error occured while verifying SSL certificate: %v", err)
	}
	return tlsConfig, nil
}

func (p *PostgreSQLProvider) DB() *gorm.DB {

	return p.db
}
