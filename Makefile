.PHONY: test
test:
	go test -v ./...

.PHONY: swagger-generate
swagger-generate:
	swag init -g /internal/app/app.go -o api/docs

.PHONY: delete-images
delete-images:
	chmod ugo+x tools/images.sh && tools/./images.sh

.PHONY: docker-build
docker-build:
	docker build -t auth-service-image .

.PHONY: run
run:
	docker compose down && docker compose up --build -d

.PHONY: initilize
initilize:
	chmod ugo+x tools/initilize.sh && tools/./initilize.sh

.PHONY: migrate
migrate:
	migrate -path migrations/ -database "postgresql://postgres:test123456@localhost:5446/authdb?sslmode=disable" -verbose up

.PHONY: stop
stop:
	docker compose down
