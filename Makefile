.PHONY: dc run test lint

dc:
	docker-compose up --remove-orphans --build -d

run:
	go build -o app cmd/audit_log/main.go && ./app

test:
	go test -race ./...

lint:
	golangci-lint run