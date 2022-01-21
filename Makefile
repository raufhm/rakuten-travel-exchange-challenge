postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e  POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root rtx_test

dropdb:
	docker exec -it postgres12 dropdb rtx_test

gotodb:
	docker exec -it postgres12 psql rtx_test

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/rtx_test?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/rtx_test?sslmode=disable" -verbose down

.PHONY: createdb dropdb postgres migrateup migratedown