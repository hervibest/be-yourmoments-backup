# Load environment variables from .env file
### === Start Services ===
start-user-svc:
	cd user-svc/cmd/web && go run main.go

start-photo-svc:
	cd photo-svc/cmd/web && go run main.go

start-upload-svc:
	cd upload-svc/cmd/web && go run main.go

start-transaction-svc:
	cd transaction-svc/cmd/web && go run main.go


include .env
# expor

.PHONY: all proto migrate-up migrate-down \
        photo-svc-migrate-up photo-svc-migrate-down photo-svc-migrate-reset \
        user-svc-migrate-up user-svc-migrate-down user-svc-migrate-reset \
        transaction-svc-migrate-up transaction-svc-migrate-down transaction-svc-migrate-reset \
        start-photo-svc start-upload-svc start-user-svc start-transaction-svc \
        mockgen-upload-svc


### === Migration Template ===
migrate-up:
	@echo "⚠️  Please set DB_URL when calling this target, e.g., make migrate-up DB_URL=$(USER_DB_URL)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	@echo "⚠️  Please set DB_URL when calling this target, e.g., make migrate-down DB_URL=$(USER_DB_URL)"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

### === User Service Migration ===
user-svc-migrate-up:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" up

user-svc-migrate-down:
	cd user-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" down

user-svc-migrate-reset:
	cd user-svc && \
	goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" down-to 0 && \
	goose -dir $(MIGRATIONS_DIR) postgres "$(USER_DB_URL)" up


### === Photo Service Migration ===
photo-svc-migrate-up:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" up

photo-svc-migrate-down:
	cd photo-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" down

photo-svc-migrate-reset:
	cd photo-svc && \
	goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" down-to 0 && \
	goose -dir $(MIGRATIONS_DIR) postgres "$(PHOTO_DB_URL)" up


### === Transaction Service Migration ===
transaction-svc-migrate-up:
	cd transaction-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" up

transaction-svc-migrate-down:
	cd transaction-svc && goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" down

transaction-svc-migrate-reset:
	cd transaction-svc && \
	goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" down-to 0 && \
	goose -dir $(MIGRATIONS_DIR) postgres "$(TRANSACTION_DB_URL)" up

### === Protobuf Compile ===
proto:
	cd pb && \
	protoc -I=. -I=pb \
	  --go_out=paths=source_relative:. \
	  --go-grpc_out=paths=source_relative:. \
	  pb/$(PROTO_FILE)

### === Mock Generation ===
mockgen-upload-svc:
	cd upload-svc/internal && \
	mockgen -source=./adapter/user_adapter.go \
	        -destination=./mocks/adapter/mock_user_adapter.go \
	        -package=mockadapter
