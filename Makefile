

all:
	make build-main

build-main:
	mkdir -p dist/ && \
        go build -o dist/mcversions ./cmd/mcversions

package:
	./scripts/package-mcversions.sh

deb:
	make clean && \
	make all && \
	make package

clean:
	rm -rf dist/ && \
	rm -rf package/


.PHONY: run
run:
	cd ./dist/ && \
	./mcversions
