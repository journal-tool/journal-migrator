# Journal-migrator

Evolving a database schema is often a scary operation for the database administrators.
Many schema changes seem safe at first, but may introduce suttle differences on how the applications
interacts with the databases, or worse: data loss. In order to safely evolve the schemas, and keep
and keep a record of all the changes across your organization databases, use Journal.

This component is a library, meant to be used to execute the schema changes
in a synchronous or asynchronous manner (i.e. without downtime).

## 🧑‍💻 Usage

This component can be run from the console, and provides the following commands:
- Migrate: to execute a schema-change, asynchronously or synchronously.
- Cleanup: to remove a schema-change leftover table, synchronously.

### Migrate command

```shell
go run cmd/migrate/main.go [ARGUMENTS]
```

List of arguments:

| Name             | Env variable  | Required | Default | Possible values            |
|------------------|---------------|----------|---------|----------------------------|
| databaseUsername | DATABASE_USER | True     | -       | *                          |
| databasePassword | DATABASE_PASS | True     | -       | *                          |
| databaseHost     | DATABASE_HOST | True     | -       | *                          |
| databasePort     | DATABASE_PORT | True     | -       | *                          |
| databaseName     | DATABASE_NAME | True     | -       | *                          |
| databaseType     | DATABASE_TYPE | True     | -       | {mysql, postgres}          |
| logLevel         | LOG_LEVEL     | False    | INFO    | {DEBUG, INFO, WARN, ERROR} |
| table            | -             | True     | -       | *                          |
| operations       | -             | True     | -       | *                          |
| strategy         | -             | False    | SYNC    | {ASYNC, SYNC}              |

The `ASYNC` strategy has two pre-requisites:
- The target table must have a Primary Key with one auto-incremented column.
- The database instance must have Foreign Key enforcement disabled.

### Cleanup command

```shell
go run cmd/cleanup/main.go [ARGUMENTS]
```

List of arguments:

| Name             | Env variable  | Required | Default | Possible values            |
|------------------|---------------|----------|---------|----------------------------|
| databaseUsername | DATABASE_USER | True     | -       | *                          |
| databasePassword | DATABASE_PASS | True     | -       | *                          |
| databaseHost     | DATABASE_HOST | True     | -       | *                          |
| databasePort     | DATABASE_PORT | True     | -       | *                          |
| databaseName     | DATABASE_NAME | True     | -       | *                          |
| databaseType     | DATABASE_TYPE | True     | -       | {mysql, postgres}          |
| logLevel         | LOG_LEVEL     | False    | INFO    | {DEBUG, INFO, WARN, ERROR} |
| table            | -             | True     | -       | *                          |


## 🔧 Development

### Dependencies
```shell
make install-dev
```

### Linting
```shell
make check
```

### Testing
```shell
make test-unit
```

To test a MySQL database:
```shell
export DATABASE_PORT=3306
export DATABASE_TYPE=mysql
export COMPOSE_FILE=mysql-8.0.yaml

make services-up && make test-integration
make services-down
```

To test a PostgreSQL database:
```shell
export DATABASE_PORT=5432
export DATABASE_TYPE=postgres
export COMPOSE_FILE=postgres-16.yaml

make services-up && make test-integration
make services-down
```

### Release
```shell
make tag
```


## Acknowledgements

This project is a Golang port of the popular [LHM][repo-lhm] Ruby project,
generalizing it to support both MySQL 8.0+ and PostgreSQL 16+ databases.


[repo-lhm]: https://github.com/soundcloud/lhm
