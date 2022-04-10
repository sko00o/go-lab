package redis

import (
	"time"
)

type Config struct {
	URL         string        `mapstructure:"url"`
	MaxIdle     int           `mapstructure:"max_idle"`
	MaxActive   int           `mapstructure:"max_active"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}
