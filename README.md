# Bookmark

## Setup(for development)

```
docker-compose up -d
docker exec -it bookmark_server /bin/bash
# run server 
go run main.go
```

```
# database login
docker exec -it bookmark_db psql -U bookmark
```

## Unittest

```
docker exec -it bookmark_server /bin/bash
./unittest.sh
```

## Down

```
docker-compose down -v
```