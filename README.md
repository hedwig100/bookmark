# Bookmark

## Setup(for development)

```
docker-compose up -d
docker exec -it bookmark_server /bin/bash
```

```
# database login
docker exec -it bookmark_db psql -U bookmark
```

## Down

```
docker-compose down -v
```