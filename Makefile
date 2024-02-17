run:
	mkdir -p ./bin/config
	mkdir -p ./bin/graphql
	go build -o ./bin/UpworkLeadgen ./cmd/UpworkLeadgen.go
	cp ./config/* ./bin/config
	cp ./internal/upwork/graphql/* ./bin/graphql
	cd ./bin && ./UpworkLeadgen

build:
	docker-compose up --build

rebuild:
	docker-compose down
	docker-compose up --build -d