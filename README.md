# service-template

## Generation api by openapi spec

```sh
oapi-codegen --config=src/api/oapi-codegen.yaml src/api/openapi.yaml    
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