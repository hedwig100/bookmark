# Bookmark

## Setup

### Backend
```

git clone https://github.com/hedwig100/bookmark && cd bookmark
docker-compose up -d

# database login
docker exec -it bookmark_db psql -U bookmark

# If you access localhost:8083, you can see an api document.

# test (latter of them failed unless you re-setup the database or 'docker-compose down -v')
docker exec bookmark_server bash ./unittest.sh
docker exec bookmark_server go test ./server/ -v -run TestIntegrate*

```

If you want to run the server locally and access to the server (ex. `https://localhost:8081/hello` ) in your browser, 
you have to set your domain and create your own ssl certificate (Chrome doesn't allow us to use self-signed certificate). 

```
# run the server 
docker exec -it bookmark_server go run main.go 
# You can access https://<your-domain>:8081/hello in the browser. 

```

## Frontend

You have to set your domain and create your own ssl certificate for the same reason above. You may use the same certificate 
as the server for simplicity. 

```
cd frontend
npm install 
npm run serve
# You can access https://<your-domain>:8084/login in the browser.
```

## Down

```
docker-compose down -v
```