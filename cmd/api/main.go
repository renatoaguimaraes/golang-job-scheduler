package main

import (
	"flag"
	"log"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/api"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
)

func main() {
	config := conf.NewConfig()
	flag.StringVar(&config.ServerAddress, "host", "localhost:8080", "host:port")
	flag.StringVar(&config.ClientCA, "ca", "cert/client-ca-cert.pem", "client ca path")
	flag.StringVar(&config.ServerCertificate, "cert", "cert/server-cert.pem", "server cert path")
	flag.StringVar(&config.ServerKey, "key", "cert/server-key.pem", "server key path")
	flag.Parse()
	if err := api.StartServer(config); err != nil {
		log.Fatalf("fail to start server, %v", err)
	}
}
