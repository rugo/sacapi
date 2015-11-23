KEYDIR=keys
KEY=$(KEYDIR)/key.pem
SNAKEOIL=tools/tls_key_params.txt
TARGET=sacapi
RM=rm -rf
.PHONY: install all clean


all: $(TARGET) $(KEY)

$(KEY): $(SNAKEOIL)
	tools/create_test_cert.sh < $< 

$(TARGET): main.go
	go get
	go build

clean:
	$(RM) $(TARGET)

mrproper: clean
	$(RM) $(KEYDIR)
