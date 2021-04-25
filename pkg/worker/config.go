package worker

import (
	"os"
)

type Config struct {
	logFolder string
}

func NewConfig() Config {
	return Config{logFolder: os.TempDir()}
}

func (c *Config) LogFolder() string {
	return c.logFolder
}
