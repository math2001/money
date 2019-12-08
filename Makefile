.PHONEY: run

run: money
	./money server

money: *.go
	go build

money-tools: tools/*.go
	go build -o money-tools ./tools

