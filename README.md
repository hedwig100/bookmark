# Bookmark

## Setup

```

git clone https://github.com/hedwig100/bookmark && cd bookmark
docker-compose up -d

# run server and you can access localhost:8081/hello
docker exec -it bookmark_server go run main.go

# database login
docker exec -it bookmark_db psql -U bookmark

# If you access localhost:8083, you can see an api document.

# test (latter of them failed unless you re-setup the database or 'docker-compose down -v')
docker exec bookmark_server bash ./unittest.sh
docker exec bookmark_server go test ./server/ -v -run TestIntegrate*

```

## Down

```
docker-compose down -v
```