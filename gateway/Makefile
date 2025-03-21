
.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: generate-coverage
generate-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

.PHONY: generate-swagger
generate-swagger:
	swag init -d "./" -g "cmd/api/main.go"

.PHONY: docker-run-all-services
docker-run-all-services:
	docker compose -f docker-compose.yml -f quickstart-standalone.yml up --build --force-recreate
