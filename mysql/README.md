# MySQL

The MySQL dialect uses the [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) driver, which [passes the compatibility test suite](https://github.com/golang/go/wiki/SQLDrivers).

By default, this MySQL dialect will parse `DATE` and `DATETIME` columns into `[]byte` or `string` types. Support for `time.Time` must be explicitly enabled by adding the `parseTime=true` parameter to the connection string.


### Testing

A valid MySQL connection string should be set on the environmental variable `SOL_TEST_MYSQL`. An example:

    user:pass@tcp(host:port)/db?parseTime=true

This variable can be given inline:

    SOL_TEST_MYSQL="user:pass@tcp(host:port)/db?parseTime=true" go test

If the environmental variable is not given, it will default to a [Travis CI](https://docs.travis-ci.com/user/database-setup/#MySQL) connection string, which will likely panic on your local system.

#### Docker

Docker hub provides an [official MySQL image](https://hub.docker.com/_/mysql/). A container can be started with:

    docker run -p 3306:3306 --name mysql -e "MYSQL_ROOT_PASSWORD=" \
    -e "MYSQL_ALLOW_EMPTY_PASSWORD=yes" \
    -e "MYSQL_DATABASE=sol_test" -d mysql:latest

Then change the host of the connection string to the default Docker machine and run tests, for example:

    SOL_TEST_MYSQL="root@tcp(<DOCKER HOST>:3306)/sol_test?parseTime=true" go test
