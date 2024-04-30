postgres:
	docker run --name postgres12 -p 5200:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root moneystack
dropdb:
	docker exec -it postgres12 dropdb moneystack

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5200/moneystack?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5200/moneystack?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY:postgres  createdb dropdb migrateup migratedown sqlc test