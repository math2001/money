.PHONY: run pwawatch build clean ocrserver

run: money
	./money

build:
	# build the service worker
	tsc -p tsconfig-sw.json
	# build the regular typescript files
	tsc
	go build

clean:
	go clean
	rm -r ./pwa/js/*

ocrserver:
	docker run -it -p 31563:8080 otiai10/ocrserver

pwawatch:
	# DOESN'T BUILD THE SERVICE WORKER
	# you can set that up yourself by doing tsc -p tsconfig-sw.json --watch
	tsc --watch

money: $(shell fd --extension go)
	go build

money-tools: tools/*.go
	go build -o money-tools ./tools


