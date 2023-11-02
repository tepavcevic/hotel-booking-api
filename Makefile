build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running api inside docker container"
	@docker run -p 8080:8080 api

test:
	@go test -v ./...