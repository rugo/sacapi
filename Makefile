.PHONY: install all clean


all: main.go
	go get
	go build

clean:
	rm -rf sacapi
