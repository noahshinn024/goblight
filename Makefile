SHELL := /bin/bash
.PHONY: all clean

base-dir := $(shell pwd)

build:
	mkdir -p out; \
		GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o goblight; \
		mv goblight $(base-dir)/out/goblight;

run:
	mkdir -p out; \
		GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o goblight; \
		mv goblight $(base-dir)/out/goblight; \
		$(base-dir)/out/goblight;


install:
	mkdir -p out; \
		GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o goblight; \
		mv goblight $(base-dir)/out/goblight; \
		chmod -R 777 $(base-dir)/out/; \
		mv $(base-dir)/out/goblight /usr/bin/goblight; \
		chown root:root /usr/bin/goblight; \
		chmod 4755 /usr/bin/goblight; 

clean:
	rm -fr out;
