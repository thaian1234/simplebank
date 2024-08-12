#!make
include app.env
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
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres12 dropdb simple_bank
migrateup:
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose up
migrateup1:
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose up 1
migratedown:
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose down
migratedown1:
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose down 1
resetdb:
	docker exec -it postgres12 dropdb simple_bank
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank
	migrate -path db/migrations -database "${DB_SOURCE}" -verbose up
sqlc:
	sqlc generate
	mockgen -package mockdb -destination db/mock/store.go github.com/thaian1234/simplebank/db/sqlc Store
test:
	go test -v -cover ./...
server: 
	go run main.go
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/thaian1234/simplebank/db/sqlc Store
.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 resetdb sqlc 