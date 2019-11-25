.PHONEY: run

run: money
	./money

money: *.go
	go build

money-tools: tools/*.go
	go build ./tools -o money-tools

