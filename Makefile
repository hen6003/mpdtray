# instantmenu - menu for instantOS
# See LICENSE file for copyright and license details.

all: build

build: 
	go build mpdtray.go

clean:
	rm -f mpdtray

install:
	cp mpdtray $(DESTDIR)$(PREFIX)/bin/
	chmod +x $(DESTDIR)$(PREFIX)/bin/mpdtray

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/mpdtray

.PHONY: all build  clean install uninstall
