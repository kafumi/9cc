9cc: $(wildcard *.go)
	go build -o $@ $^

.PHONY: test
test: 9cc
	./test.sh
