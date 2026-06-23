package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port        int           `mapstructure:"port"`
	Environment string        `mapstructure:"environment"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`

	AuthServiceAddr        string `mapstructure:"auth_service_addr"`
	UserServiceAddr        string `mapstructure:"user_service_addr"`
	DocumentIngestionAddr string `mapstructure:"document_ingestion_addr"`
	SearchServiceAddr     string `mapstructure:"search_service_addr"`
	NotificationServiceAddr string `mapstructure:"notification_service_addr"`

	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPassword string `mapstructure:"redis_password"`
	RedisDB       int    `mapstructure:"redis_db"`

	KafkaBrokers      []string `mapstructure:"kafka_brokers"`
	KafkaClientID     string   `mapstructure:"kafka_client_id"`

	RateLimit         int     `mapstructure:"rate_limit"`
	RateLimitBurst    int     `mapstructure:"rate_limit_burst"`

	JWTSecret         string  `mapstructure:"jwt_secret"`

	OTLPEndpoint      string  `mapstructure:"otel_endpoint"`

	CorsAllowedOrigins []string `mapstructure:"cors_allowed_origins"`
}

func Load() *Config {
	viper.SetDefault("port", 8081)
	viper.SetDefault("environment", "development")
	viper.SetDefault("read_timeout", "30s")
	viper.SetDefault("write_timeout", "30s")
	viper.SetDefault("idle_timeout", "120s")
	viper.SetDefault("redis_db", 0)
	viper.SetDefault("kafka_client_id", "api-gateway")
	viper.SetDefault("rate_limit", 100)
	viper.SetDefault("rate_limit_burst", 200)
	viper.SetDefault("otel_endpoint", "localhost:4318")
	viper.SetDefault("cors_allowed_origins", []string{"http://localhost:3000"})

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GW")

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}
	return cfg
}
