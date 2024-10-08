# Hozirgi direktoriyani saqlaydi
CURRENT_DIR := $(shell pwd)
DATABASE_URL="postgres://postgres:0509@localhost:5432/artisanconnect_aut?sslmode=disable"

# Serverni ishga tushiradi
runserver:
	@go run cmd/router/server.go

# Xizmatni ishga tushiradi
runservice:
	@go run cmd/service/service.go

# Protobuf fayllarni generatsiya qiladi
gen-proto:
	@./scripts/gen-proto.sh $(CURRENT_DIR)

# Go modullarni tozalaydi
tidy:
	@go mod tidy

# Yangi migratsiya yaratadi
mig-create:
	@if [ -z "$(name)" ]; then \
		read -p "Migratsiya nomini kiriting: " name; \
	fi; \
	migrate create -ext sql -dir migrations -seq $$name

# Mavjud migratsiyalarni qo'llaydi
mig-up:
	@migrate -database "$(DATABASE_URL)" -path migrations up

# Oxirgi migratsiyalarni bekor qiladi
mig-down:
	@migrate -database "$(DATABASE_URL)" -path migrations down

# Migratsiya versiyasini majburiy ravishda o'zgartiradi
mig-force:
	@if [ -z "$(version)" ]; then \
		read -p "Migratsiya versiyasini kiriting: " version; \
	fi; \
	migrate -database "$(DATABASE_URL)" -path migrations force $$version

# `gen-proto.sh` skriptiga ruxsat beradi
permission:
	chmod +x scripts/gen-proto.sh

# Testlarni ishga tushiradi
test:
	@go test ./storage/postgres

swager:
	swag init -g api/api.go --output api/docs/

	