package providers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/SZabrodskii/microservices/libs/providers/config"
	"github.com/go-redis/redis/v8"
	"go.uber.org/fx"
	"os"
)

type RedisProvider struct {
	client *redis.Client
}

type RedisProviderOptions struct {
	Options redis.Options
}

func NewRedisProvider(cfg config.RedisConfig) (*RedisProvider, error) {
	var opts redis.Options
	if cfg.GetOptions() != nil {
		opts = *cfg.GetOptions()
	}

	opts.Addr = cfg.GetDSN()
	opts.Password = cfg.GetPassword()
	opts.DB = cfg.GetDB()

	var tlsCfg *tls.Config
	if cfg.GetTLSConfig() != nil {
		cert, err := tls.LoadX509KeyPair(cfg.GetTLSConfig().GetCertificate(), cfg.GetTLSConfig().GetKey())
		if err != nil {
			return nil, fmt.Errorf("%w: could not get TLS Certificates", err)
		}

		caCert, err := os.ReadFile(cfg.GetTLSConfig().GetRootCertificate())
		if err != nil {
			return nil, fmt.Errorf("%w: could not get root certificate", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsCfg = &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}

	opts.TLSConfig = tlsCfg

	client := redis.NewClient(&opts)
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("error occurred while connecting to the Redis server: %w", err)
	}

	provider := &RedisProvider{
		client: client,
	}

	return provider, nil
}

func (rprvd *RedisProvider) getTLSConfig(certificatePath, rootCertificatePath string) (*tls.Config, error) {
	certificate, err := os.ReadFile(certificatePath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading SSL certificate: %w", err)
	}

	rootCertificate, err := os.ReadFile(rootCertificatePath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading SSL root certificate: %w", err)
	}

	parsedCertificate, err := x509.ParseCertificate(certificate)
	if err != nil {
		return nil, fmt.Errorf("error occured while parsing certificate: %w", err)
	}

	rootCertificates, err := x509.ParseCertificates(rootCertificate)
	if err != nil {
		return nil, fmt.Errorf("error parsing root certificate: %w", err)
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
		return nil, fmt.Errorf("error occured while verifying SSL certificate: %w", err)
	}
	return tlsConfig, nil
}

func (rprvd *RedisProvider) Client() *redis.Client {
	return rprvd.client
}

//=========================================== End of Experimenting Code ===============================================

func NewRedisSentinelProvider(config config.RedisSentinelConfig, tlsConfig *tls.Config) (*RedisProvider, error) {
	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       config.GetMasterName(),
		SentinelAddrs:    config.GetHosts(),
		Password:         config.GetPassword(),
		SentinelPassword: config.GetSentinelPassword(),
		DB:               config.GetDB(),
	})

	prvd := &RedisProvider{
		client: client,
	}

	return prvd, nil

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
