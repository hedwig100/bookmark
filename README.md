# Bookmark

## Setup(only database now)

```
docker-compose up -d
docker exec -it bookmark_db psql -U bookmark
```

## Down

```
docker-compose down -v
```