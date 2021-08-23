.PHONY: build run clean

build:
	go build -o alpha ./cmd/*.go

run:
	go build -o alpha ./cmd/*.go
	./alpha

clean:
	rm alpha
	rm *.result.csv
