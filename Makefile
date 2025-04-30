PHOTO_DB_URL=postgres://postgres:postgres@localhost:5432/photo_svc?sslmode=disable&TimeZone=Asia/Jakarta
USER_DB_URL=postgres://postgres:postgres@localhost:5432/user_svc?sslmode=disable&TimeZone=Asia/Jakarta
TRANSACTION_DB_URL=postgres://postgres:postgres@localhost:5432/transaction_svc?sslmode=disable&TimeZone=Asia/Jakarta

MIGRATIONS_DIR=db/migrations
PROTO_DIR=user-svc/internal/pb
PROTO_FILE=user.proto

.PHONY: migrate-down proto

start-photo-svc:
	cd photo-svc/cmd/web && go run main.go

start-upload-svc:
	cd upload-svc/cmd/web && go run main.go

start-user-svc:
	cd user-svc/cmd/web && go run main.go

start-transaction-svc:
	cd transaction-svc/cmd/web && go run main.go

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down


photo-svc-migrate-down:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" down

photo-svc-migrate-up:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" up

photo-svc-migrate-reset:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" down-to 0 &&  goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" up 


user-svc-migrate-down:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" down

user-svc-migrate-up:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" up

user-svc-migrate-reset:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" down-to 0 &&  goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" up 


transaction-svc-migrate-down:
	cd transaction-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" down

transaction-svc-migrate-up:
	cd transaction-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" up

transaction-svc-migrate-reset:
	cd transaction-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" down-to 0 &&  goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" up
	
proto:
	cd pb && protoc --go_out=. --go-grpc_out=. $(PROTO_FILE)

mockgen-upload-svc:
	cd upload-svc/internal && mockgen -source=./adapter/user_adapter.go -destination=./mocks/adapter/mock_user_adapter.go -package=mockadapter