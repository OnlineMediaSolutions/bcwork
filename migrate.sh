#!/bin/bash
docker stop pg

docker run --rm   --name pg  -d -p 5433:5432 -e POSTGRES_PASSWORD=postgres  postgres || exit $?

export PGPASSWORD=postgres
pg_isready -h localhost -U postgres -p 5433

while [ $? -ne 0 ]
do
   echo "waiting for db init...."
   sleep 1
   pg_isready -h localhost -U postgres -p 5433
done

echo "db ready"

PGPASSWORD=postgres

psql -h localhost -U postgres -p 5433 < init.sql > /dev/null
migrate -source file://migrations -database postgres://postgres:postgres@localhost:5433/bcdb-dev?sslmode=disable goto 12

echo "sqlboiler wipe"
sqlboiler psql --wipe
#docker stop pg
