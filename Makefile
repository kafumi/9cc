.PHONY: build
build:
	docker run --rm -v $(shell pwd):/usr/src/myapp -w /usr/src/myapp golang:1.12 make -f docker.mk 9cc

.PHONY: test
test:
	docker run --rm -v $(shell pwd):/usr/src/myapp -w /usr/src/myapp golang:1.12 make -f docker.mk test

.PHONY: clean
clean:
	rm -f 9cc tmp* test/*.o