SHELL := /bin/bash

.PHONY: build
build:
	go build -o system-test -v ./...

.PHONY: unit
unit:
	go test -cover ./. ./internal/...

.PHONY: tag-dockerfile
tag-dockerfile:
	docker tag huautla:latest huautla:lkg
	
.PHONY: postgres
postgres:
	docker-compose up --build --remove-orphans -d postgres

.PHONY: inspect
inspect:
	docker-compose exec -it postgres /bin/sh -c "psql -hlocalhost -Upostgres huautla"

.PHONY: system-test
system-test: postgres
	docker-compose exec \
		-ePOSTGRES_HOST=localhost \
		-ePOSTGRES_PORT=5432 \
		-ePOSTGRES_USER=postgres \
		-ePOSTGRES_PASSWORD=root \
		postgres \
    /bin/sh -c "cd /huautla && ./bin/install-system-test.sh"
	# -docker-compose up --build --remove-orphans system-test
	docker tag huautla:latest huautla:lkg
	make docker-down

vet:
	go vet ./...

fmt:
	go fmt ./...

.PHONY: docker-down
docker-down:
	docker-compose down --remove-orphans
