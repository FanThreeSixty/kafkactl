# Makefile provides a simple interface for building and testing Kafkactl
build:
	cd commands; go mod tidy
	cd commands; go mod vendor
	go mod tidy
	go mod vendor
	go build -o kafkactl-darwn .

test:
	echo "test"