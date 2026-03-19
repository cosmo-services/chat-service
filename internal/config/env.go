package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type Env struct {
	AppEnv string `mapstructure:"APP_ENV"`
	Port   string `mapstructure:"PORT"`

	PGHost string `mapstructure:"CHAT_PG_HOST"`
	PGPort string `mapstructure:"CHAT_PG_PORT"`
	PGUser string `mapstructure:"CHAT_PG_USER"`
	PGPass string `mapstructure:"CHAT_PG_PASS"`
	PGName string `mapstructure:"CHAT_PG_NAME"`

	NatsHost string `mapstructure:"NATS_HOST"`
	NatsPort string `mapstructure:"NATS_PORT"`

	AppDomain string `mapstructure:"APP_DOMAIN"`

	MigrationPath string `mapstructure:"MIGRATION_PATH"`

	JwtSecret string `mapstructure:"JWT_SECRET"`

	AllowedOrigins []string `mapstructure:"ALLOWED_ORIGINS"`

	GrpcPort                 string `mapstructure:"GRPC_PORT"`
	SocialServiceGrpcAddress string `mapstructure:"GRPC_SOCIAL_ADR"`
}

func NewEnv() Env {
	env := Env{}

	_, err := os.Stat(".env")
	useEnvFile := !os.IsNotExist(err)

	if useEnvFile {
		viper.SetConfigType("env")
		viper.SetConfigName(".env")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal("Can't read the .env file: ", err)
		}

		err = viper.Unmarshal(&env)
		if err != nil {
			log.Fatal("Environment can't be loaded: ", err)
		}
	} else {
		env.bindEnv()
	}

	if env.AppEnv != "production" {
		log.Println("The App is running in development env")
	}

	return env
}

func (e *Env) bindEnv() {
	e.AppEnv = os.Getenv("APP_ENV")
	e.Port = os.Getenv("PORT")

	e.PGHost = os.Getenv("CHAT_PG_HOST")
	e.PGPort = os.Getenv("CHAT_PG_PORT")
	e.PGUser = os.Getenv("CHAT_PG_USER")
	e.PGPass = os.Getenv("CHAT_PG_PASS")
	e.PGName = os.Getenv("CHAT_PG_NAME")

	e.NatsHost = os.Getenv("NATS_HOST")
	e.NatsPort = os.Getenv("NATS_PORT")

	e.AppDomain = os.Getenv("APP_DOMAIN")

	e.JwtSecret = os.Getenv("JWT_SECRET")

	e.MigrationPath = os.Getenv("MIGRATION_PATH")

	e.GrpcPort = os.Getenv("GRPC_PORT")
	e.SocialServiceGrpcAddress = os.Getenv("GRPC_SOCIAL_ADR")

	if val := os.Getenv("ALLOWED_ORIGINS"); val != "" {
		e.AllowedOrigins = strings.Split(val, ",")
	}
}

var Module = fx.Options(
	fx.Provide(NewEnv),
)
