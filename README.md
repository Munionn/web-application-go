# Project webapplication

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

**Live reload** (like nodemon for Node.js – rebuilds on file changes):
```bash
make watch
```
Uses [Air](https://github.com/air-verse/air). On first run, Air is installed via `go install` if needed. Edit any `.go` file to trigger a rebuild and restart.

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
