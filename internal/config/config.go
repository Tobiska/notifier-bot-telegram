package config

import (
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var configInstance = &Config{}
var once sync.Once

type Config struct {
	Telegram
	Logger
	Database
}

type Telegram struct {
	BotFatherToken string        `env:"TELEGRAM_BOT_FATHER_TOKEN"`
	TimeoutUpdates int           `env:"TELEGRAM_TIMEOUT_UPDATES"`
	Debug          bool          `env:"TELEGRAM_BOT_DEBUG"`
	RetryCount     int           `env:"TELEGRAM_RETRY_COUNT"`
	RetryTimeout   time.Duration `env:"TELEGRAM_RETRY_DURATION"`
}

type Logger struct {
	Lever string `env:"LOGGER_LEVEL"`
}

type Database struct {
	Dsn string `env:"DATABASE_DSN"`
}

func ReadConfig() (*Config, error) {
	var configErr error
	once.Do(func() {
		configErr = cleanenv.ReadEnv(configInstance)
	})
	return configInstance, configErr
}
