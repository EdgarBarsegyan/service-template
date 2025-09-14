# service-template

Пример сервиса для быстрой разработки и интеграции

## Generation api by openapi spec

```sh
cd src/
oapi-codegen --config=api/oapi-codegen.yaml api/openapi.yaml    
```

## Generate migration
```sh
migrate create -ext sql -dir src/internal/persistence/infrastructure/migrations -seq <MigrationName>
```

## Run db into docker

- create and run container
```sh
sudo docker run --name postgres-test -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=test -e POSTGRES_DB=testdb -p 5432:5432 -d postgres:latest
```

- start container
```sh
 sudo docker start postgres-test
```