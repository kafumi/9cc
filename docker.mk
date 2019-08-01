9cc: 9cc.go
	go build $^

.PHONY: test
test: 9cc
	./test.sh
