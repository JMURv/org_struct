package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	LogLvl      string `env:"LOG_LEVEL"    envDefault:"debug"`
	Mode        string `env:"MODE"         envDefault:"dev"`
	ServiceName string `env:"SERVICE_NAME" envDefault:"org_struct"`
	Server      ServerConfig
	DB          dbConfig
	Jaeger      jaegerConfig
}

type ServerConfig struct {
	Scheme   string `env:"SERVER_SCHEME"             envDefault:"http"`
	Domain   string `env:"SERVER_DOMAIN"             envDefault:"localhost"`
	Port     int    `env:"SERVER_HTTP_PORT,required"`
	GRPCPort int    `env:"SERVER_GRPC_PORT"          envDefault:"50050"`
	PromPort int    `env:"SERVER_PROM_PORT"          envDefault:"8085"`
}
type dbConfig struct {
	Host     string `env:"POSTGRES_HOST"     envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT"     envDefault:"5432"`
	User     string `env:"POSTGRES_USER"     envDefault:"app_owner"`
	Password string `env:"POSTGRES_PASSWORD" envDefault:"password"`
	Database string `env:"POSTGRES_DB"       envDefault:"app_db"`
}

type jaegerConfig struct {
	Sampler struct {
		Type  string  `env:"JAEGER_SAMPLER_TYPE" envDefault:"const"`
		Param float64 `env:"JAEGER_SAMPLER_PARAM" envDefault:"1"`
	} `yaml:"sampler"`
	Reporter struct {
		LogSpans           bool   `env:"JAEGER_REPORTER_LOGSPANS" envDefault:"true"`
		LocalAgentHostPort string `env:"JAEGER_REPORTER_LOCALAGENT" envDefault:"localhost:6831"`
		CollectorEndpoint  string `env:"JAEGER_REPORTER_COLLECTOR" envDefault:"http://localhost:14268/api/traces"`
	} `yaml:"reporter"`
}

func MustLoad() Config {
	const defaultPath = "config/.env"
	if err := godotenv.Load(defaultPath); err != nil {
		if !os.IsNotExist(err) {
			panic("failed to load .env file: " + err.Error())
		}
		log.Println("No .env file found, using system environment variables")
	} else {
		log.Println("Loaded environment variables from: " + defaultPath)
	}

	conf := Config{}
	if err := env.Parse(&conf); err != nil {
		panic("failed to parse environment variables: " + err.Error())
	}

	log.Println("Load configuration from environment")

	return conf
}
