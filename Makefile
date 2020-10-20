.PHONY:test

build:
	go build  -o ./bin/mongotest ./mongo/*.go
	
test:
	go test -v ./test/

clean:
	rm bin/*
