APPLICATION_NAME=toggl
ADMIN_USER ?= "admin"
ADMIN_PASS ?= "admin"
MARIADB_IMAGE ?= mariadb
MARIADB_VERSION ?= 10.5.9
LOCAL_DB_IMAGE = $(MARIADB_IMAGE):$(MARIADB_VERSION)
LOCAL_DB_USER ?= $(APPLICATION_NAME)
LOCAL_DB_PASS ?= $(APPLICATION_NAME)
LOCAL_DB_NAME ?= $(APPLICATION_NAME)
LOCAL_DB_PORT ?= 3307

LOCAL_PORT ?= 8089


.PHONY: help clean build  run
.DEFAULT_GOAL := help

export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit


clean: ## clean the backend build results
	go mod download github.com/BurntSushi/toml
	go clean
	rm -rf bin

build-backend: ## build the backend for the target platforms
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/toggl-linux-amd64 cmd/main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/toggl-darwin-amd64 cmd/main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/toggl-darwin-arm64 cmd/main.go

build: build-backend ## build the frontend and the backend

stop:
	ps aux|grep bin/toggl|grep -v
run: build## run the broker locally
		BIND_ADDRESS=127.0.0.1:$(LOCAL_PORT) \
		MARIADB_USER=$(LOCAL_DB_USER) \
		MARIADB_PASS=$(LOCAL_DB_PASS) \
		MARIADB_DB=$(LOCAL_DB_NAME) \
		MARIADB_HOST=127.0.0.1 \
		MARIADB_PORT=$(LOCAL_DB_PORT) \
		ADMIN_USER=${ADMIN_USER} \
		ADMIN_PASS=${ADMIN_PASS} \
		bin/toggl-$$(go env GOOS)-$$(go env GOARCH)


local-db-start: local-db-stop ## start local database for running APPLICATION tests
	docker pull -q $(LOCAL_DB_IMAGE)
	docker run --name $(APPLICATION_NAME)-db -d \
		-e MYSQL_DATABASE=$(LOCAL_DB_NAME) \
		-e MYSQL_USER=$(LOCAL_DB_USER) \
		-e MYSQL_PASSWORD=$(LOCAL_DB_PASS) \
		-e MYSQL_RANDOM_ROOT_PASSWORD=true \
		-p 127.0.0.1:$(LOCAL_DB_PORT):3306 \
	    $(LOCAL_DB_IMAGE)

local-db-stop: ## stop local database
	-docker rm -f $(APPLICATION_NAME)-db

local-db-connect: ## connect to local database
	docker exec -it $(APPLICATION_NAME)-db mysql --user=$(LOCAL_DB_USER) --password=$(LOCAL_DB_PASS) $(LOCAL_DB_NAME)
lint:
	gofmt
