SOURCE=$(wildcard *.go)
TARGET=sacregister

.PHONY=all install clean

all: $(TARGET)

$(TARGET): $(SOURCE)
	go build

install: $(SOURCE)
	go install

clean:
	$(RM) $(TARGET) 
