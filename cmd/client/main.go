package main

import (
	"os"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/client/command"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
)

func main() {
	// TODO move to a configuration file
	config := conf.NewConfig()
	config.ServerAddress = "localhost:8080"
	config.ServerCA = "cert/server-ca-cert.pem"
	config.ClientCertificate = "cert/client-cert.pem"
	config.ClientKey = "cert/client-key.pem"
	command.Execute(config, os.Args[1:])
}
