package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type key string

const KeyUUID key = key("uuid")
const KeyMetrics key = key("metrics")

type Config struct {
	Service   Service
	Postgres  Postgres
	Kafka     Kafka
	Metrics   Metrics
	Platform  Platform
	Optionhub Optionhub
}

type Service struct {
	Port string `env:"USER_SERVICE_PORT"`
}

type Postgres struct {
	User     string `env:"USER_SERVICE_POSTGRES_USER"`
	Password string `env:"USER_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"USER_SERVICE_POSTGRES_DB"`
	Host     string `env:"USER_SERVICE_POSTGRES_HOST"`
	Port     string `env:"USER_SERVICE_POSTGRES_PORT"`
}

type Kafka struct {
	Server          string `env:"KAFKA_SERVER"`
	FriendsRegister string `env:"USER_FRIENDS_REGISTER"`
	SetNewAvatar    string `env:"AVATAR_SET_NEW"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

type Optionhub struct {
	Host string `env:"OPTIONHUB_SERVICE_HOST"`
	Port string `env:"OPTIONHUB_SERVICE_PORT"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
