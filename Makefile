build:
	go build -o bin/mongotest
test:
	go test -v
clean:
	rm bin/*
