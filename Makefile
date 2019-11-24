.PHONEY: run

run: money
	./money

money: *.go
	go build

