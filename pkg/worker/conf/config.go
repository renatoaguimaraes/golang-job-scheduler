package conf

import "os"

// Config worker configuration.
type Config struct {
	// LogFolder stores all job logs
	LogFolder string
	// LogChunckSize size in bytes for each log chunck read from log file
	LogChunckSize int

	ServerAddress string

	ServerCA          string
	ServerCertificate string
	ServerKey         string

	ClientCA          string
	ClientCertificate string
	ClientKey         string
}

func NewConfig() Config {
	return Config{
		LogFolder:     os.TempDir(),
		LogChunckSize: 1024,
	}
}
