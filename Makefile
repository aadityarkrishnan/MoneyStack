postgres:
  docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
  docker exec -it postgres12 createdb --username=root --owner=root moneystack

dropdb:
  docker exec -it postgres12 dropdb moneystack

migrateup:
  migrate -path db/migration -database "$(DB_SOURCE)" -verbose up

migratedown:
  migrate -path db/migration -database "$(DB_SOURCE)" -verbose down

migrateup_one:
  migrate -path db/migration -database "$(DB_SOURCE)" -verbose up 1

migratedown_one:
  migrate -path db/migration -database "$(DB_SOURCE)" -verbose down 1

sqlc:
  sqlc generate

test:
  go test -v -cover ./...

server:
  go run main.go

mock:
  mockgen -package mockdb -destination db/mock/store.go github/aadityarkrishnan/MoneyStack/db/sqlc Store

.PHONY:postgres  createdb dropdb migrateup migratedown migrateup_one migratedown_one  sqlc test server mock
