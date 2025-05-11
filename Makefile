test:
	docker-compose -f docker-compose-test.yaml up -d
	go test -v ./...
	docker-compose -f docker-compose-test.yaml down --volumes
.PHONY: test

up:
	docker-compose up --build
.PHONY: up

down:
	docker-compose down
.PHONY: down

downvolumes:
	docker-compose down --volumes
.PHONY: down