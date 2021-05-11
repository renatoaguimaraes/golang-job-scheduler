# Job Scheduler

[![CI](https://github.com/renatoaguimaraes/job-scheduler/actions/workflows/ci.yml/badge.svg?branch=library)](https://github.com/renatoaguimaraes/job-scheduler/actions/workflows/ci.yml)

## Summary

Prototype job worker service that provides an API to run arbitrary Linux processes.

## Overview

### Library

*   Worker library with methods to start/stop/query status and get an output of a running job.

### API

*   Use GRPC for API to start/stop/get status of a running process;
*   Add streaming log output of a running job process; 
*   Use mTLS and verify client certificate; 
*   Set up a strong set of cipher suites for TLS and a good crypto setup for certificates;
*   Authentication and Authorization.

### Client	

*   Client command should be able to connect to worker service and schedule several jobs. 
*   The client should be able to query the results of the job execution and fetch the logs. 
*   The client should be able to stream the logs.

## Security

### Transport

The Transport Layer Security (TLS), version 1.3, provides privacy and data integrity in secure communication between the client and server.

### Authentication

The authentication will be provided by mTLS. The following assets will be generated, using the openssl v1.1.1k, to support the authorization schema:

* Server CA private key and self-signed certificate
* Server private key and certificate signing request (CSR)
* Server signed certificate, based on Server CA private key and Server CSR
* Client CA private key and self-signed certificate
* Client private key and certificate signing request (CSR)
* Client signed certificate, based on Client CA private key and Client CSR

The authentication process checks the certificate signature, finding a CA certificate with a subject field that matches the issuer field of the target certificate, once the proper authority certificate is found, the validator checks the signature on the target certificate using the public key in the CA certificate. If the signature check fails, the certificate is invalid and the connection will not be established. Both client and server execute the same process to validate each other. Intermediate certificates won't be used.

### Authorization
The user roles will be added into the client certificate as an extension, so the gRPC server interceptors will read and check the roles to authorize the user, the roles available are reader and writer, the writer will be authorized do start, stop, query and stream operations and the reader will be authorized just to query and stream operations available on the API. A memory map will keep the name of the gRPC method and the respective roles authorized to execute it.
The X.509 v3 extensions will be used to add the user role to the certificate. For that, the extension attribute roleOid 1.2.840.10070.8.1 = ASN1:UTF8String, must be requested in the Certificate Signing Request (CSR), when the user certificate is created, the user roles must be informed when the CA signs the the CSR. The information after UTF8String: is encoded inside of the x509 certificate under the given OID.

#### gRPC interceptors
* UnaryInterceptor
* StreamInterceptor

#### Certificates
* X.509
* Signature Algorithm: sha256WithRSAEncryption
* Public Key Algorithm: rsaEncryption
* RSA Public-Key: (4096 bit)
* roleOid 1.2.840.10070.8.1 = ASN1:UTF8String (for the client certificate)

## Out of scope

*   Database to persist the worker state;
*   Log rotation, purge data policy, distributed file system;

## Design Proposal 

![Architecture](assets/architecture.jpg)

## Test

```sh
$ make test
```
## Build and run API

```sh
$ make api
go build -o ./bin/worker-api cmd/api/main.go
```

```sh
$ ./bin/worker-api
```

## Build and run Client

```sh
$ make client
go build -o ./bin/worker-client cmd/client/main.go
```

```sh
$ ./bin/worker-client start "bash" "-c" "while true; do date; sleep 1; done"
Job 9a8cb077-22da-488f-98b4-d2fb51ba4fc9 is started
```

```sh
$ ./bin/worker-client query 9a8cb077-22da-488f-98b4-d2fb51ba4fc9
Pid: 1494556 Exit code: 0 Exited: false
```

```sh
$ ./bin/worker-client stream 9a8cb077-22da-488f-98b4-d2fb51ba4fc9
Sun 02 May 2021 05:54:29 PM -03
Sun 02 May 2021 05:54:30 PM -03
Sun 02 May 2021 05:54:31 PM -03
Sun 02 May 2021 05:54:32 PM -03
Sun 02 May 2021 05:54:33 PM -03
Sun 02 May 2021 05:54:34 PM -03
Sun 02 May 2021 05:54:35 PM -03
Sun 02 May 2021 05:54:36 PM -03
Sun 02 May 2021 05:54:37 PM -03
```

```sh
./bin/worker-client stop 9a8cb077-22da-488f-98b4-d2fb51ba4fc9
Job 79d95817-7228-4c36-8054-6c29513841b4 has been stopped
```
