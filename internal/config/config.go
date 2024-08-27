package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

var configInstance *Config
var once sync.Once

type Config struct {
	Telegram
}

type Telegram struct {
	BotFatherToken string `env:"TELEGRAM_BOT_FATHER_TOKEN"`
	TimeoutUpdates int    `env:"TELEGRAM_TIMEOUT_UPDATES"`
	Debug          bool   `env:"TELEGRAM_BOT_DEBUG"`
}

func ReadConfig() (*Config, error) {
	var configErr error
	once.Do(func() {
		configErr = cleanenv.ReadEnv(configInstance)
	})
	return configInstance, configErr
}
