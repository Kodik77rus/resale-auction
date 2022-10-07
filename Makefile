all: build
build:
	go build -o ./bin/auction ./cmd/auction
clean:
	rm -f ./bin/auction
run:
	./bin/auction -p 8080 -d 0.0.0.0:6767,0.0.0.0:4354
dev:
	make build && make run