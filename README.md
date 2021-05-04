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

## Out of scope

*   Database to persist the worker state;
*   Log rotation, purge data policy, distributed file system;
*   Authorization server;
*   Cluster setup;
*   Configuration files.

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
