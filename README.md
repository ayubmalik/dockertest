# dockertestspike

A spike/example application for runnig golang integration test in parallel using [dockertest](https://github.com/ory/dockertest) and [PostgreSQL](https://www.postgresql.org/).

## Usage

You can run tests with or _without_ using dockertest. This is useful for comparing what dependencies you need on your system when running tests (e.g. postgresql, docker etc)


### Run tests _without_ Dockertest

You will need to have [postgresql](https://www.postgresql.org/) up and running locally before running `go test`. The database named `dockertest` must exist. As the tests run against the _same_ DB instance can't really run them in parallel, hence the `test.parallel` param is set to 1.

```
go test -test.parallel 1
```

### Run tests _with_ Dockertest

You only need to have docker up and running. The first time the tests are run, the docker [postgresql image](https://hub.docker.com/_/postgres?tab=tags&page=1&ordering=last_updated&name=13.2-alpine) will be pulled, so may be slower.

```
go test -dt test.parallel 4
```
