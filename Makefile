.PHONY:

run:
	go run main.go

build:	## Build backend Docker image
	docker build . \
		-t shop-api:latest

docker.run:
	docker run -d \
	-p 8080:8080 \
	--name shop-api shop-api

##########
#  sqlc  #
##########

include Makefile.env

sqlc.gen:
	sqlc -f ./database/sqlc/sqlc.yaml generate

migrate.new:
	goose sqlite -dir=${MYSQL_SQL_PATH} . create . sql

migrate.validate:
	goose -dir=${MYSQL_SQL_PATH} validate

migrate.up:
	goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(DB_NAME)?charset=utf8&parseTime=True&loc=Local" -dir=${MYSQL_SQL_PATH} up

migrate.down:
	goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(DB_NAME)?charset=utf8&parseTime=True&loc=Local" -dir=${MYSQL_SQL_PATH} reset 