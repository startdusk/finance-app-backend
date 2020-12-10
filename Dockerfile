FROM golang:1.15.3 as builder
ENV DATA_DIRECTORY /go/src/github.com/startdusk/finance-app-backend
WORKDIR $DATA_DIRECTORY
ARG APP_VERSION
ARG CGO_ENABLED=0

## I added this part to avoid redownloading all libraries on each build, now it caches all libraries and use them if go.mod wasn't changed(我添加这个部分是为了避免在每个构建中重新加载所有库，现在它缓存所有库并在go.mod没有改变的情况下使用它们)
COPY go.mod .
COPY go.sum .
RUN go env -w GOPROXY=https://goproxy.cn
RUN go mod download
## those 4 lines

COPY . .
RUN go build -ldflags="-X github.com/startdusk/finance-app-backend/internal/config.Version=$APP_VERSION" github.com/startdusk/finance-app-backend/cmd/server

FROM alpine:3.10
ENV DATA_DIRECTORY /go/src/github.com/startdusk/finance-app-backend
RUN apk add --update --no-cache \ 
	ca-certificates
COPY ./internal/database/migrations ${DATA_DIRECTORY}/internal/database/migrations
COPY --from=builder ${DATA_DIRECTORY}/server /finance-app-backend

ENTRYPOINT ["/finance-app-backend"]
