package conf

import "os"

// Config worker configuration.
type Config struct {
	// LogFolder stores all job logs
	LogFolder string
	// LogChunckSize size in bytes for each log chunck read from log file
	LogChunckSize int
}

func NewConfig() Config {
	return Config{
		LogFolder:     os.TempDir(),
		LogChunckSize: 1024,
	}
}
