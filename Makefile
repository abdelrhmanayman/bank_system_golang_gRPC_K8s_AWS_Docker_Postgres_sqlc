createdb:
	docker exec -it postgres12 createdb --username=root --owner=root bank_db
	
dropdb:
	docker exec  -it postgres12 dropdb bank_db

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/bank_db?sslmode=disable" -verbose down 1

test: 
	go test -v -cover ./...

postgres: 
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

sqlc:
	sqlc generate

resetDB:
	make migratedown && make migrateup

server:
	go run index.go

serverDev:
	nodemon --exec go run index.go --signal SIGTERM

mockDB: 
	mockgen -package mockdb -destination db/mock/store.go banksystem/db/sqlc Store
	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test resetDB server serverDev mockDB migratedown1
