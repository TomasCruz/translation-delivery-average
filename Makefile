.PHONY: all
all: clean build

.PHONY: clean
clean:
	go clean
	rm -f ./translation-delivery-average

.PHONY: build
build:
	CGO_ENABLED=0 go build -o ./translation-delivery-average

.PHONY: fresh
fresh: db all

.PHONY: run
run:
	./translation-delivery-average --input_file events.txt --window_size 10

.PHONY: db_start
db_start:
	docker start delivery_data

.PHONY: db
db: db_restart
db_restart:
	docker rm -f delivery_data &> /dev/null
	sleep 1
	docker run --name delivery_data -e POSTGRES_USER=toma -e POSTGRES_PASSWORD=pswd -p 5372:5432 -d postgres
	docker start delivery_data
	sleep 1

psql:
	docker run -it --rm --link delivery_data:postgres postgres psql postgresql://toma:pswd@delivery_data
