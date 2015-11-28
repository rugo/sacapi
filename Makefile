MODULES=sacregister saccalendar
BINDIR=modules/service
RM=rm -rf
.PHONY: all clean install keys

all clean install:
	for dir in $(MODULES); do \
	    $(MAKE) -C $(BINDIR)/$$dir $@; \
	done

keys:
	tools/create_test_cert.sh < tools/tls_key_params.txt
