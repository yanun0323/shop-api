# build stage
FROM golang:1.23-alpine AS build

ADD . /go/build
WORKDIR /go/build

ADD ./config/config.yaml /go/build/
ADD go.mod go.sum /go/build/

RUN go mod download

# install gcc
RUN apk add build-base

RUN go build -o shop-api main.go

# final stage
FROM alpine:3.20

# install timezone data
RUN apk add --no-cache tzdata
ENV TZ Asia/Taipei
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /go/build/shop-api /var/application/shop-api
COPY --from=build /go/build/config.yaml /var/application/config.yaml

EXPOSE 8080

WORKDIR /var/application
CMD [ "./shop-api" ]