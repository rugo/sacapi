SOURCE=$(wildcard *.go)
TARGET=register

.PHONY=all clean

all: $(TARGET)

$(TARGET): $(SOURCE)
	go build

clean:
	$(RM) $(TARGET) 
