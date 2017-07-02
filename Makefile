.PHONY: test

test_full:
	go test -v -race -benchmem -bench=.

bench:
	go test -v -bench=. -cpuprofile=cpu.prof

benchmem:
		go test -v -benchmem -bench=. -cpuprofile=cpu.prof

test:
	go test -v

deps:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u honnef.co/go/tools/cmd/gosimple
	go get -u honnef.co/go/tools/cmd/unused

check:
	staticcheck $$(go list ./... | grep -v /vendor/)
	gosimple $$(go list ./... | grep -v /vendor/)
	unused $$(go list ./... | grep -v /vendor/)