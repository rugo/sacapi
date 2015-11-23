KEY=keys/key.pem
SNAKEOIL=tools/tls_key_params.txt

.PHONY: install all clean


all: main.go $(KEY)
	go get
	go build

$(KEY): $(SNAKEOIL)
	tools/create_test_cert.sh < $< 

clean:
	rm -rf sacapi
