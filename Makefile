MODULES=sacregister saccalendar
BINDIR=modules/service
RM=rm -rf
.PHONY: all clean keys

all:
	 for dir in $(MODULES); do \
	  $(MAKE) -C $(BINDIR)/$$dir; \
     done

install:
	for dir in $(MODULES); do \
    	  $(MAKE) -C $(BINDIR)/$$dir $@; \
    done

clean:
	for dir in $(MODULES); do \
    	  $(MAKE) -C $(BINDIR)/$$dir $@; \
    done
