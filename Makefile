# Makefile provides a simple interface for building and testing Kafkactl
build:
	cd commands; rm -rf vendor
	cd commands; go mod tidy
	cd commands; go mod vendor
	rm -rf vendor
	go mod tidy
	go mod vendor
	go build -o kafkactl-darwin .

test:
	echo "test"