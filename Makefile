createdb:
	docker exec -it postgres-bank createdb --username=root --owner=root bank
dropdb:
	docker exec -it postgres-bank dropdb bank
run-postgres:
	docker run --name postgres-bank -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e 'POSTGRES_DB=bank' -p 5432:5432 -d postgres:15.2-alpine 
run-admin:
	docker run -p 80:80 -e 'PGADMIN_DEFAULT_EMAIL=admin@admin.com' -e 'PGADMIN_DEFAULT_PASSWORD=admin'-d dpage/pgadmin4
migrateup:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" --verbose up
migratedown:
	migrate --path db/migrations --database "postgresql://root:secret@localhost:5432/bank?sslmode=disable" --verbose down 
sqlc:
	sqlc generate
test:
	go test --v --cover ./...

.PHONY: test sqlc
