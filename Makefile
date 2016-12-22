.PHONY: test

test_full:
	go test -v -race -benchmem -bench=.

bench:
	go test -v -bench=.

benchmem:
		go test -v -benchmem -bench=.

test:
	go test -v
