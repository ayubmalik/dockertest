# dockertestspike

A spike/example application for runnig golang integration test in parallel using [dockertest](https://github.com/ory/dockertest) and [PostgreSQL](https://www.postgresql.org/).

## Usage

You only need to have docker up and running. The first time the tests are run, the docker [postgresql image](https://hub.docker.com/_/postgres?tab=tags&page=1&ordering=last_updated&name=13.2-alpine) will be pulled, so may be slower.

```
go test -dt test.parallel 4
```
