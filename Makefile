.PHONY: go_fmt go_vet go_tidy go_get go_test

go_fmt:
	docker compose run --rm api go fmt ./...

go_vet:
	docker compose run --rm api go vet ./...

go_tidy:
	docker compose run --rm api go mod tidy

go_get:
	docker compose run --rm api go get ${pkg}

go_test:
	docker compose up -d test-db
	docker compose run --rm api go test -v -cover ./...
	docker compose stop test-db
