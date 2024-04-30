# MoneyStack
MoneyStack is an app to track the expense between their friends &amp; family

DB Diagram for the MS: https://dbdiagram.io/d/MoneyStack-658856e789dea627997c92bc

To pull an postgres12: docker pull postgres:12-alpine

Check if image is availble use docker images

To create the container, run the image.

docker run --name postgres12 -p 5200:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

To get the running container docker ps

To exec a command inside the container
docker exec -it postgres12 psql -U root

To get the logs of the container
docker logs 5dc958d273c8

DB Migration
https://github.com/golang-migrate/migrate

To install 
brew install golang-migrate

mkdir -p db/migration

migrate create -ext sql -dir db/migration -seq init_schema 

This will create  the up & down sql files

To create the DB in docker

docker exec -it postgres12 createdb --username=root --owner=root moneystack
docker exec -it postgres12 psql -U root simple_bank

For performign the migration

Put the postgres query from the db diagram  to migrate up file, and drop commands in the down file

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5200/moneystack?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5200/moneystack?sslmode=disable" -verbose down

SQL Quering

SQLC will be using for quering purpose

brew install kyleconroy/sqlc/sqlc

sqlc init 

 To initialize the sqlc use inside the root of the project and it will create the yml file

 Add these to sqlc.yaml

 version: "2"
packages:
  - name: "db"
    path: "db/sqlc"
    queries: "/db/query/"
    schema: "/db/migration/"
    engine: "postgresql"
    emit_json_tags: true
    emit_prepared_queries: false
    emit_interface: false
    emit_exact_table_names: false


`go mod init github/aadityarkrishnan/MoneyStack` is a command in Go programming language used to initialize a new Go module. A Go module is a collection of related Go packages that are versioned together.

go mod tidy


In order to connect the go to postgres, need to install a driver

go get github.com/lib/pq (#https://github.com/lib/pq)

In the go.mod you might as indirect references, 

_ "github.com/lib/pq" // The underscore before the import will keep it even if it is not used in the code.

go mod tidy // Can be used to keeping the mod file clean

Testify is an package to perform the testing in Go

to install Testify,

go get github.com/stretchr/testify


Query to get the locked query

`SELECT blocked_locks.pid     AS blocked_pid,
         blocked_activity.usename  AS blocked_user,
         blocking_locks.pid     AS blocking_pid,
         blocking_activity.usename AS blocking_user,
         blocked_activity.query    AS blocked_statement,
         blocking_activity.query   AS current_statement_in_blocking_process,
         blocked_activity.application_name AS blocked_application,
         blocking_activity.application_name AS blocking_application
   FROM  pg_catalog.pg_locks         blocked_locks
    JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
    JOIN pg_catalog.pg_locks         blocking_locks 
        ON blocking_locks.locktype = blocked_locks.locktype
        AND blocking_locks.DATABASE IS NOT DISTINCT FROM blocked_locks.DATABASE
        AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
        AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
        AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
        AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
        AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
        AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
        AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
        AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
        AND blocking_locks.pid != blocked_locks.pid JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
   WHERE NOT blocked_locks.GRANTED;`

SELECT a.datname,
          a.application_name,
         l.relation::regclass,
         l.transactionid,
         l.mode,
         l.locktype,
         l.GRANTED,
         a.usename,
         a.query,
         a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
ORDER BY a.query_start;

In order to deploy the application we use github actions
for that create .github/workflows/ci.yaml in the root folder 