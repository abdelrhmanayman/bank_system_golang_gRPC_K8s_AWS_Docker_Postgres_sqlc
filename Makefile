createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bank_db
	
dropdb:
	docker exec  -it postgres12 dropdb bank_db

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose down

test: 
	go test -v -cover ./...

postgres: 
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

sqlc:
	sqlc generate

resetDB:
	make migratedown && make migrateup
	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test resetDB
