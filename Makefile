postgres:
		docker run --name postgres14 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:14-alpine

createdb:
	  docker exec -it postgres14 createdb --username=root --owner=root simple_bank

dropdb:
	  docker exec -it postgres14 dropdb simple_bank

migrateup:
		migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateuplast:
		migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
		migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedownlast:
		migrate -path db/migration -database "postgresql://root:root@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
		sqlc generate

server:
		go run main.go

mock:
		mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/radugaf/simplebank/db/sqlc Store

test:
	go test -v -cover ./...
	
.PHONY: postgres createdb dropdb migrateup migrateuplast migratedown migratedownlast sqlc server mock test
