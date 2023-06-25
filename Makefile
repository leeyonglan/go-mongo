.PHONY:test build

.DEFAULT_GOAL := build

-include ./build/config.ini

build:
	./build/build.sh $(type)

test:
	./bin/teaapp-dev

clean:
	rm bin/*

sync:
	./build/rsync.sh $(type)

help:
	@echo "make - complie the source code"
	@echo "make sync -- rysnc the binary file to test server"
