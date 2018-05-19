run:
	go run main/main.go

test:
	go clean -testcache
	go test -v ./...

update:
	git pull
	docker-compose down
	docker-compose up --build -d
	sleep 2
	docker-compose logs
