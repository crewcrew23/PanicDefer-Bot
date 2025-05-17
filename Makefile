.PHONY: vet test run-server run-handler run-worker docker-up docker-down migrate

vet:
	go vet ./...

test:
	go test -race ./...

run-server:
	go run  .\cmd\server\main.go --config .\config\local.yaml

run-handler:
	go run .\cmd\workerHandler\main.go --config .\config\local.yaml

run-worker:
	go run .\cmd\workerHealthzChecker\main.go --config .\config\local.yaml

migrate:
	go run .\cmd\migrator\main.go --db-url postgres://user:pass@host:port/dbname

docker-up:
	docker-compose up --build --detach

docker-down:
	docker-compose down