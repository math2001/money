.PHONEY: run dev

run: money
	./money server

dev:
	tsc --watch

money: *.go
	go build

money-tools: tools/*.go
	go build -o money-tools ./tools

