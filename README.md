# Golang Rest API

## How To Use This Repo

- **Entity-based**: The resources available should represent the domain model. Each resource should have the CRUD methods implemented (even if not all available to API consumers). In our template, we have a single resource defined (users.go). However other resources could be easily added by copying the template and changing the logic of the service layer.
- **Properly Abstracted**: The Transport, service, and data layers are all cleanly abstracted from one another. This makes it easy to make apply updates to the API endpoints
- **Consistent**: It's important that consumers of a service have guaranteed consistency across the entire range of API endpoints and methods. In this service, responses are consistently formatted whether successfully returning a JSON object or responding with an error code. All the service's methods use shared response (http.go) and error (errors.go) handler functions to ensure consistency.
- **Tested**: We believe that a blend of unit and integration testing is important for ensuring that the service maintains its contract with consumers. The service repo therefore contains a collection of unit and integration tests for the various layers of the service.
- **Explorable**: It is important for developers to be able to play with an endpoint in order to understand it. We have provided Postman collections for testing out the REST endpoints exposed by the service. That's why there is a `Bootstrap Users` collection that can be run using the `Run collection` tool in Postman that will create 100 users to test the search endpoint with.

## Getting Started

### Prerequisites

- Go 1.18 (should still be backwards compatible with earlier versions)

### Running locally

1. From root of the repo
2. Run `docker-compose up` will start the dependencies and server on port 8080

### Running via docker

1. From root of the repo
2. Run `docker-compose up` will start the dependencies and server on port 8080

### Postman

The collections will need an environment setup with `scheme`, `port` and `host` variables setup with values of `http`, `8080` and `localhost` respectively.

### Run tests

Some of the integration tests use docker to spin up dependencies on demand (ie a postgres db) so just be aware that docker is needed to run the tests.

1. From root of the repo
2. Run `go test ./...`
