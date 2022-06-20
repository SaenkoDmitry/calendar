# calendar

## Local run

#### build project:
```bash
make build
```

#### run project (without db):
```bash
make run
```

#### run calendar by one command with local database (inside docker containers):
```bash
docker-compose up --build
```

### ATTENTION!
```
when service started 'goose' utility will apply migrations from 'calendar/internal/db/migrations' including '20220621003129_generate_data.sql'
which contains postgres sql scripts for generating random mock db data for your convinience
so first start might be a bit slow
```

## Migrations

#### example of rollback to previous migration from actual
```bash
goose -dir internal/db/migrations postgres "user=postgres password=postgres dbname=calendar sslmode=disable" down
```

## Swagger

#### to open swagger page go by link
```
http://localhost:8080/swagger/index.html
```

#### to format swagger comments
```bash
swag fmt
```

#### to update swagger files
```bash
swag init -g cmd/main.go
```

## Tests

#### run tests
```bash
make test
```

#### run tests with coverage
```bash
go test -coverpkg=./... ./...
```
