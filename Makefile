.PHONEY: run pwawatch build clean

run: money
	./money server

build:
	# build the service worker
	tsc -p tsconfig-sw.json
	# build the regular typescript files
	tsc
	go build

clean:
	go clean
	rm -r ./pwa/js/*

pwawatch:
	# DOESN'T BUILD THE SERVICE WORKER
	# you can set that up yourself by doing tsc -p tsconfig-sw.json --watch
	tsc --watch

money: $(fd --extension go)
	go build

money-tools: tools/*.go
	go build -o money-tools ./tools


