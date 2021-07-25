# Go Infrastructure
[![Go Report Card](https://goreportcard.com/badge/github.com/janmbaco/go-infrastructure)](https://goreportcard.com/report/github.com/janmbaco/go-infrastructure)

This is an infrastructure project in go that serves the Go-ReverseProxy-SSL and Saprocate projects, it also aims to be a common base for projects in go:

### Package: github.com/janmbaco/go-infrastructure/logs
- It provides a Log service which allows writing to a log file from a logging level (Trace, Info, Warning, Error, Fatal)

### Package: github.com/janmbaco/go-infrastructure/errorhandler
- It provides a series of utilities for handling errors in Go that can be thrown or caught. If the log is set below fatal, all errors are logged.

### Package: github.com/janmbaco/go-infrastructure/event
- It provides a Service for subscribe and publish events.

### Package: github.com/janmbaco/go-infrastructure/disk
- It provides tools for writing and deleting files on disk, it also provides a service that listens for changes to a disk file.

###  Package: github.com/janmbaco/go-infrastructure/config
- It provides a configuration interface and an implementation for a file configuration.

### Package: github.com/janmbaco/go-infrastructure/server
- It provides a service that starts an http or gRpc server that automatically restarts when the configuration changes.

### Package: github.com/janmbaco/go-infrastructure/crypto
- It provides a service to encrypt and decrypt bytes.
