VERSION=$(shell git rev-parse --short HEAD) # 从git中获取上一次提交的git log的短hash信息作为version传递到容器里面的应用中

build-dev:
	docker-compose build --build-arg APP_VERSION=$(VERSION)

up-dev:
	docker-compose up server
