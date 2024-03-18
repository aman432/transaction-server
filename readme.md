# Transaction-server
A user transaction recording server
#### Relevant Specs and Docs:
- Postman Collection: [Link](https://documenter.getpostman.com/view/14037306/2sA2xpR8ns)

#### Local Setup (Docker)

- Run `make docker-build` to build api and migrations docker image.
  - Run `make docker-build-api` to build api docker image.
  - Run `make docker-build-migration` to build migration docker image.
  
- Run `make dev-docker-up` to start the server using docker compose. The server should be running at localhost:9040
- `curl --location --request POST 'http://localhost:9040/health/check' \
  --header 'Content-Type: application/json' \
  --data-raw '{ }'`

#### Local Setup (without Docker, Recommended)
- Run `make pre-build` this will run below-mentioned commands
    - Run `make deps` to get all the build dependencies
    - Run `make mock-gen` to generate all mock files to use in running unit test cases.
    - Run `go-build-migration` this command will build the migration binary.
    - Run `go-build-api` this command will build the api binary.
- Start the services with below commands.
    - Run `make up-migration` to run migrations. Uses mysql and expects `prizmo` db created.
    - Run `make go-run-api` to start the server. The server should be running at localhost:9040
- Run `make test` to run all the test cases.
- Run `make test-coverage` to run all the test cases and generate coverage report.
- Run `make swagger` to generate swagger documentation.
   - Run `make swagger-serve` to serve the swagger documentation at localhost:55863/docs

#### Database Setup (for without Docker)

- Install mysql/postgres and create a database `prizmo` with user `root` and password `root`.
- Kindly change config files based on your configurations.

```toml
[db]
    debug                     = false
    [db.ConnectionConfig]
        dialect               = "postgres/mysql"
        protocol              = "tcp"
        url                   = "localhost"
        port                  =  5432
        username              = "root"
        password              = "root"
        sslMode               = "disable"
        name                  = "prizmo"
    [db.ConnectionPoolConfig]
        maxOpenConnections    = 5
        maxIdleConnections    = 5
        connectionMaxLifetime = 0
```
