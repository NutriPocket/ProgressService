test:
	docker compose -f docker-compose-test.yaml up -d
	sleep 15
	cd src && go test -v ./... && cd ..
	docker compose -f docker-compose-test.yaml down --volumes
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