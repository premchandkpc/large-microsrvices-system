package config

import "github.com/spf13/viper"

type Config struct {
	Port        int      `mapstructure:"port"`
	Environment string   `mapstructure:"environment"`
	KafkaBrokers []string `mapstructure:"kafka_brokers"`
	RedisAddr    string   `mapstructure:"redis_addr"`

	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	SMTPFrom     string `mapstructure:"smtp_from"`

	PushEnabled  bool   `mapstructure:"push_enabled"`
	FirebaseKey  string `mapstructure:"firebase_key"`

	SlackToken   string `mapstructure:"slack_token"`
	SlackChannel string `mapstructure:"slack_channel"`
}

func Load() *Config {
	viper.SetDefault("port", 8087)
	viper.SetDefault("environment", "development")
	viper.SetDefault("smtp_port", 587)
	viper.SetDefault("smtp_from", "noreply@platform.example.com")
	viper.SetDefault("push_enabled", false)

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NS")

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}
	return cfg
}
