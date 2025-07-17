.PHONY: test

PLUGINS=echo

test:
	go test -cover ./...

plugins: $(PLUGINS)
$(PLUGINS):
	go get -t ./...
	go build -buildmode=plugin -o ../build/$@.so ./plugin/$@
	go test -cover ./plugin/$@
