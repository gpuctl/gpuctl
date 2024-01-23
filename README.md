# gpuctl - A GPU control room

## Running tests

Running tests with `go test ./...` requires access to a postgres installation
and role with permission to create tables in a database.

The default tests database is configured to connect to a database called
`gpuctl-tests-db` owned by a role `gpuctl`.

This can be created by running the following command as a user/role with
administrative access to postgres (typically `postgres`):

```
createuser gpuctl
createdb -O gpuctl gpuctl-tests-db
```
