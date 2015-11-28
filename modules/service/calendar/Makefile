SOURCE=$(wildcard *.go)
TARGET=calendar

.PHONY=all clean

all: $(TARGET)

$(TARGET): $(SOURCE)
	go build

clean:
	$(RM) $(TARGET)
