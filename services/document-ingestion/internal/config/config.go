package config

import "github.com/spf13/viper"

type Config struct {
	Port        int    `mapstructure:"port"`
	Environment string `mapstructure:"environment"`

	StorageType string `mapstructure:"storage_type"`
	S3Endpoint  string `mapstructure:"s3_endpoint"`
	S3Region    string `mapstructure:"s3_region"`
	S3AccessKey string `mapstructure:"s3_access_key"`
	S3SecretKey string `mapstructure:"s3_secret_key"`
	S3Bucket    string `mapstructure:"s3_bucket"`

	KafkaBrokers []string `mapstructure:"kafka_brokers"`

	RedisAddr string `mapstructure:"redis_addr"`

	MaxUploadSize int64 `mapstructure:"max_upload_size"`
	AllowedTypes  []string `mapstructure:"allowed_types"`
}

func Load() *Config {
	viper.SetDefault("port", 8084)
	viper.SetDefault("environment", "development")
	viper.SetDefault("storage_type", "s3")
	viper.SetDefault("s3_region", "us-east-1")
	viper.SetDefault("s3_bucket", "documents")
	viper.SetDefault("max_upload_size", 104857600)
	viper.SetDefault("allowed_types", []string{
		"application/pdf",
		"image/jpeg",
		"image/png",
		"text/plain",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"text/csv",
		"application/json",
	})

	viper.AutomaticEnv()
	viper.SetEnvPrefix("DI")

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}
	return cfg
}
