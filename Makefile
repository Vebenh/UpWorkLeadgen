run:
	mkdir -p ./bin/config
	go build -o ./bin/UpworkLeadgen ./cmd/UpworkLeadgen.go
	cp ./config/*.yaml ./bin/config
	cd ./bin && ./UpworkLeadgen

build:
	docker-compose up --build

rebuild:
	docker-compose down
	docker-compose up --build -d