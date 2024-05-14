.PHONY: build
build:
	go build ./cmd/...

.PHONY: run
run:
	go run ./cmd/...

.PHONY: generate-pb
generate-pb:
	./bin/buf generate

.PHONY: lint
lint:
	./bin/golangci-lint run \
	--config=.golangci.pipeline.yaml \
	./...

.PHONY: create-migration
create-migration:
	cd ./migrations && ../bin/goose create new_migration sql

.PHONY: migrate
migrate:
	cd ./migrations && ../bin/goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=postgres sslmode=disable" up
