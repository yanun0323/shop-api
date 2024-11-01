.PHONY:

run:
	go run main.go

test:
	go test -count=1 -cover ./...

############
#  docker  #
############

docker.build:	
	docker build . \
		-t shop-api:latest

docker.run:
	docker rm -f shop-api;\
	docker run -d \
	--network host \
	--name shop-api shop-api 

docker.stop:
	docker stop shop-api

####################
#  docker-compose  #
####################

compose.up:
	docker-compose -p shop-api up

compose.down:
	docker-compose down

database.up:
	docker-compose -f docker-compose-database.yaml -p database up -d

database.down:
	docker-compose -f docker-compose-database.yaml down


##########
#  sqlc  #
##########

include Makefile.env

install:
	brew install sqlc goose docker-compose

sqlc.gen:
	sqlc -f ./database/sqlc/sqlc.yaml generate


###############
#  migration  #
###############

migrate.new:
	goose sqlite -dir=${MYSQL_SQL_PATH} . create . sql

migrate.validate:
	goose -dir=${MYSQL_SQL_PATH} validate

migrate.up:
	goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(DB_NAME)?charset=utf8&parseTime=True&loc=Local" -dir=${MYSQL_SQL_PATH} up

migrate.down:
	goose mysql "${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(DB_NAME)?charset=utf8&parseTime=True&loc=Local" -dir=${MYSQL_SQL_PATH} reset 