#!/bin/bash

go run db/main.go -t -n 1000
psql -U bookmark -f db/delete_db.sql -d bookmark
psql -U bookmark -f db/init/create_db.sql -d bookmark
time psql -U bookmark -f db/init/sample.sql -d bookmark > out.log
time psql -U bookmark -f db/query.sql -d bookmark > out.log