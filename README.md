# Job Scheduler

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



