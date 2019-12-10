.PHONEY: dev build clean

run: money
	./money server

build:
	tsc
	go build

clean:
	go clean
	rm -r ./pwa/js/*

money: $(fd --extension go)
	go build

money-tools: tools/*.go
	go build -o money-tools ./tools


