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


To work with Rest API, Need to use Go Web Framework.
For that we will use Gin Frameworks

install using  `go get -u github.com/gin-gonic/gin`


Next we need to use Golang viper to setup the configuration such as DB connection based on the environment which would help in reusing them without redundancy.

install using `go get github.com/spf13/viper`


Now we need to mock db for testing various scenarios

go install github.com/golang/mock/mockgen@v1.6.0
go get github.com/golang/mock/mockgen/model

Then need to its path

 ls -l  ~/go/bin                                 
total 140464
-rwxr-xr-x@ 1 aadityaradhakrishnan  staff  16677122 Mar 13 09:28 dlv
-rwxr-xr-x@ 1 aadityaradhakrishnan  staff   3516722 Mar 13 09:27 go-outline
-rwxr-xr-x@ 1 aadityaradhakrishnan  staff  29442530 Mar 13 09:29 gopls
-rwxr-xr-x  1 aadityaradhakrishnan  staff   9528098 May  2 08:49 mockgen
-rwxr-xr-x@ 1 aadityaradhakrishnan  staff  12738770 Mar 13 09:29 staticcheck


which mockgen // if you are getting error that means, path is not added.

echo $SHELL // to identity which shell you're using

if your shell is /bin/zsh

then edit vi ~/.zprofile

and ADD export PATH=$PATH:~/go/bin

source ~/.zprofile

then again check which mockgen


To create the mock file


mockgen -package <package_name> -destination <destination path> <module> <Interface> 

mockgen -package mockdb -destination db/mock/store.go github/aadityarkrishnan/MoneyStack/db/sqlc Store

This will store.go in mock folder.


To create a new migration use this
migrate create -ext sql -dir db/migration -seq add_users
/Users/aadityaradhakrishnan/Coding/MoneyStack/db/migration/000002_add_users.up.sql
/Users/aadityaradhakrishnan/Coding/MoneyStack/db/migration/000002_add_users.down.sql


then add the new migration into the new file for both up/down

THen perform the migrate down to clean up existing data to avoid anamolies

the add the migrate up and down for last one in the Make file.

Now add the sql query for user table.

Then run make sqlc , and write the test for user 

then run make mock


For JWT implementation, need to install,

go get github.com/google/uuid

go get github.com/dgrijalva/jwt-go

For PASETO add
go get github.com/o1egl/paseto

After implementing the user authentication

Now need to deploy the application using docker & kub8

For the first need to create DockerFile in the root

Then in that docker file specfiy where base(go lang), working dir , copy the file from local to docker , take build using RUN command and finally the exec command

After this 
docker build -t moneystack:latest .

We divide the build stage and run build to reduce the size of docker image

Need to rebuild once made changes in DockerFile

 docker build -t moneystack:latest .

 To rebuild without cache
 docker build --no-cache -t moneystack:latest .


 TO create the docker container use this
 docker run --name moneystack -p 8080:8080 -e GIN_MODE=release  moneystack:latest

 We need to add DB details 

 docker inspect postgres12 | grep '"IPAddress"' | head -n 1 | awk -F '"' '{print $4}' // to get postgres IP address

172.17.0.1 is the bridge address  of postgres | 172.17.0.2 is IP address
 docker run --name moneystack -p 8080:8080 -e GIN_MODE=release -e  DBSOURCE="postgresql://root:root@172.17.0.2:5200/moneystack?sslmode=disable" moneystack:latest


 aadityaradhakrishnan@Aadityas-MacBook-Air MoneyStack % docker network ls
NETWORK ID     NAME              DRIVER    SCOPE
19e585af8b4e   bridge            bridge    local
0e2ce7b76cb7   docker_gwbridge   bridge    local
8d3f0541a150   host              host      local
d2b791bcb5bc   none              null      local
aadityaradhakrishnan@Aadityas-MacBook-Air MoneyStack % 

to get further details on network

docker network inspect <network name>

In order to create  network

docker network create <network_name>

In order to connect network with container

docker network connect <network_name> <container>

Now to create the moneystack  container with network


docker run --name moneystack --network moneystack_nw -p 8080:8080 \
  -e GIN_MODE=release \
  -e DBSOURCE="postgresql://root:root@postgres12:5432/moneystack?sslmode=disable" \
  moneystack:latest



To enter into the container

docker exec -it e9a614a1def8 /bin/sh


Once created the start.sh need to changes its permission to chmod +x start.sh

Need to download
https://github.com/eficode/wait-for/releases/tag/v2.2.4 (download the latest version frm this link)
 mv ~/Downloads/wait-for ./wait-for.sh
 chmod +x wait-for.sh



 docker-compose up // to run 

docker-compose up --build // to rebuild

docker-compose run --service-ports api /bin/sh // to manual run for testing


 Now we need to push the docker repo to ECR

 For that Go AWS Console and create a private ECR repository

 Then rename old github yaml to .test and change the name

 Now create a deploy.yaml to add ecr push codes
 
Now search for Github MarketPlace and look for AWS ECR and get its config code

In order to get access to ECR, need to create a new user in IAM and need to create a Permission Group with this policy AmazonEC2ContainerRegistryFullAccess

Then generate the access key and secret key and keep them on Github Secret



