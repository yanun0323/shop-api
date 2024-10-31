.PHONY:

run:
	go run main.go

migrate.new:
	goose sqlite -dir=./database/migration . create . sql

migrate.up:
	goose mysql "test:test@/shop?parseTime=true" -dir=./database/migration up

migrate.down:
	goose mysql "test:test@/shop?parseTime=true" -dir=./database/migration down

migrate.reset:
	goose mysql "test:test@/shop?parseTime=true" -dir=./database/migration reset