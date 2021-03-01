# dockertestspike

A spike/example application for runnig golang integration test in parallel using [dockertest](https://github.com/ory/dockertest) and PostgreSQL.

## Usage

You can run tests with or _without_ using dockertest. This is useful for comparing what dependencies you need on your system when running tests (e.g. postgresql, docker etc)


### Run tests _without_ Dockertest

You will need to have postgresql up and running locally before running `go test`. The database named `dockertest` must exist.

```
go test -test.parallel 1
```

### Run tests _with_ Dockertest

You only need to have docker up and running.

```
go test -dt test.parallel 4
```
