createdb:
	docker exec -it postgres-bank createdb --username=root --owner=root bank
dropdb:
	docker exec -it postgres-bank dropdb bank
create-network:
	docker network create bank-network
create-postgres:
	docker run --name postgres-bank --network bank-system -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e 'POSTGRES_DB=bank' -p 5432:5432 -d postgres:15.2-alpine 
create-postgres-admin:
	docker run --name postgres-bank-admin --network bank-system  -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin'-p 80:80 -d dpage/pgadmin4
migrateup:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" --verbose up
migrateup1:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" --verbose up 1
migratedown:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" --verbose down 
migratedown1:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" --verbose down 1
sqlc:
	sqlc generate
test:
	go test --v --cover ./...
server:
	go run main.go
mock:
	mockgen --package mockdb --destination db/mock/store.go github.com/crackz/simple-bank/db/sqlc Store

.PHONY: test sqlc server mock
