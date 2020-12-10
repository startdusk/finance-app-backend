FROM golang:1.15.3 as builder
ENV DATA_DIRECTORY /go/src/github.com/startdusk/finance-app-backend
WORKDIR $DATA_DIRECTORY
ARG APP_VERSION
ARG CGO_ENABLED=0
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn
RUN go build -ldflags="-X github.com/startdusk/finance-app-backend/internal/config.Version=$APP_VERSION" github.com/startdusk/finance-app-backend/cmd/server

FROM alpine:3.10
ENV DATA_DIRECTORY /go/src/github.com/startdusk/finance-app-backend
RUN apk add --update --no-cache \ 
	ca-certificates
COPY ./internal/database/migrations ${DATA_DIRECTORY}/internal/database/migrations
COPY --from=builder ${DATA_DIRECTORY}/server /finance-app-backend

ENTRYPOINT ["/finance-app-backend"]
