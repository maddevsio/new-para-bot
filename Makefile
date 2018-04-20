run:
	go run main/main.go

test:
	go test -v ./...

update:
	git pull
	docker-compose down
	docker-compose up --build -d
