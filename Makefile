VERSION=$(shell git rev-parse --short HEAD) # 从git中获取上一次提交的git log的短hash信息作为version传递到容器里面的应用中

build-dev: fmt-code
	docker-compose build --build-arg APP_VERSION=$(VERSION)

up-dev:
	docker-compose up server

fmt-code:
	go fmt ./...

# docker build 前先清空之前生成错误的image签TAG为<none>的容器
clear-none-docker-image:
	docker rm $(docker ps -a -q) 
	docker rmi  $(docker images -f "dangling=true" -q)
