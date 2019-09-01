SRCS=$(wildcard *.go)
TEST_SRCS=$(wildcard test/*.c)
TEST_OBJS=$(TEST_SRCS:.c=.o)

9cc: $(SRCS)
	go build -o $@ $^

.PHONY: test
test: 9cc $(TEST_OBJS)
	./test.sh
