PREFIX ?= /usr/local/bin

.PHONY: build install uninstall clean

build:
	go build -o jitzu .
	ln -sf jitzu jz

install: build
	install -m 755 jitzu $(PREFIX)/jitzu
	ln -sf $(PREFIX)/jitzu $(PREFIX)/jz

uninstall:
	rm -f $(PREFIX)/jitzu $(PREFIX)/jz

clean:
	rm -f jitzu jz
