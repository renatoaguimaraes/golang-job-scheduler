package api

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/renatoaguimaraes/job-scheduler/internal/worker/proto"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker"
	"github.com/renatoaguimaraes/job-scheduler/pkg/worker/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func loadTLSCredentials(conf conf.Config) (credentials.TransportCredentials, error) {
	// Load certificate of the CA who signed client's certificate
	pemClientCA, err := ioutil.ReadFile(conf.ClientCA)
	if err != nil {
		return nil, err
	}
	// Certification pool to append client CA's certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemClientCA) {
		return nil, fmt.Errorf("failed to add client CA's certificate")
	}
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(conf.ServerCertificate, conf.ServerKey)
	if err != nil {
		return nil, err
	}
	// Configure credentials to require and verify the client cert
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}
	return credentials.NewTLS(config), nil
}

func StartGRPCServer(conf conf.Config) error {
	lis, err := net.Listen("tcp", conf.ServerAddress)
	if err != nil {
		return err
	}
	tlsCredentials, err := loadTLSCredentials(conf)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(UnaryAuthInterceptor),
		grpc.StreamInterceptor(StreamAuthInterceptor),
	)
	proto.RegisterWorkerServiceServer(grpcServer, &workerServer{
		Worker: worker.NewWorker(conf),
	})
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}
