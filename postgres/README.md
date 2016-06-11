# PostGres

The PostGres dialect uses the [github.com/lib/pq](https://github.com/lib/pq) driver, which [passes the compatibility test suite](https://github.com/golang/go/wiki/SQLDrivers).


### Testing

A valid PostGres connection string should be set on the environmental variable `SOL_TEST_POSTGRES`. An example:

    SOL_TEST_POSTGRES="host=localhost port=5432 dbname=sol_test user=postgres password=secret sslmode=disable" go test

If the environmental variable is empty, the test will default to a [Travis CI](https://docs.travis-ci.com/user/database-setup/#PostgreSQL) connection string, which will likely panic on your local system.

The testing database must have the `uuid-ossp` enabled.

#### Docker

Docker hub provides an [official PostGres image](https://hub.docker.com/_/postgres/). A container can be started with:

    docker run -p 5432:5432 --name postgres \
    -e "POSTGRES_PASSWORD=" \
    -e "POSTGRES_DB=sol_test" -d postgres:latest

Then change the host of the connection string to the default Docker machine and run tests, for example:

    SOL_TEST_POSTGRES="host=<DOCKER HOST> port=5432 dbname=sol_test user=postgres sslmode=disable" go test
