package config

import (
	"os"

	"github.com/FlyKarlik/auth-service/internal/errs"
)

type Config struct {
	ServiceName string
	ServerHost  string
	DatabaseURL string
	JaegerHost  string
	LogLevel    string
	JWTSecret   string
}

func New() (*Config, error) {
	cfg := &Config{
		ServiceName: os.Getenv("SERVICE_NAME"),
		ServerHost:  os.Getenv("SERVER_HOST"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JaegerHost:  os.Getenv("JAEGER_HOST"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) validate() error {

	configMap := map[string]struct {
		value string
		err   error
	}{
		"ServiceName": {c.ServiceName, errs.ErrServiceNameNotConfigured},
		"ServerHost":  {c.ServerHost, errs.ErrServerHostNotConfigured},
		"DatabaseURL": {c.DatabaseURL, errs.ErrDatabaseUrlNotConfigured},
		"JaegerHost":  {c.JaegerHost, errs.ErrJaegerHostNotConfigured},
		"LogLevel":    {c.JaegerHost, errs.ErrJaegerHostNotConfigured},
		"JwtSecret":   {c.JWTSecret, errs.ErrJwtSecreNotConfigured},
	}

	for _, val := range configMap {
		if len(val.value) == 0 {
			return val.err
		}
	}

	return nil
}
