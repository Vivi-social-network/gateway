# build
FROM golang:1.23 AS build

WORKDIR /app

COPY .. .

RUN go env -w GOPRIVATE=github.com/Vivi-social-network

RUN go mod tidy
RUN go build -a -o gateway ./cmd

# run
FROM golang:1.23-alpine

ENV CONFIG_PATH=./configs/dev.yml

WORKDIR /app

COPY /configs configs
COPY --from=build /app/gateway .

CMD ./gateway -cfg $CONFIG_PATH
