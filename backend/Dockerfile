FROM golang:1.18-bullseye

EXPOSE 8081
WORKDIR /home/server
COPY . /home/server/
RUN apt-get update && \
    apt-get install openssl && \
    go mod tidy && \
    go build
ENTRYPOINT ["bash","entrypoint.sh"]