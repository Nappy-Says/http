postgres:
	docker run --name postgres10 -p 5432:5432 -e POSTGRES_USER=app -e POSTGRES_PASSWORD=pass -d postgres:10

createdb:
	docker exec -it postgres10 createdb --username=app --owner=app db

dropdb:
	docker exec -it postgres10 dropdb db

migrateup:
	migrate -path db/migration -database "postgresql://app:pass@localhost:5432/db?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://app:pass@localhost:5432/db?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown