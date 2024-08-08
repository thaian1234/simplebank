#!make

# Kiểm tra xem DB_SOURCE có được định nghĩa không
ifeq ($(origin DB_SOURCE),undefined)
    # Nếu không, thử load từ app.env
    ifneq (,$(wildcard app.env))
        include app.env
        export $(shell sed 's/=.*//' app.env)
    endif
endif

# Đảm bảo DB_SOURCE có giá trị
ifndef DB_SOURCE
    $(error DB_SOURCE is not set. Please set it in app.env or as an environment variable)
endif

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose up
migratedown:
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose down
resetdb:
	docker exec -it postgres12 dropdb simple_bank
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose up
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server: 
	go run main.go
.PHONY: postgres createdb dropdb migrateup migratedown resetdb sqlc 