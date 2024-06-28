prefix = 
DESTDIR = /usr/local
##
## pos-relay Makefile help
##
## Usage: make [target]
##
## Target includes: build
##
.PHONY: help
help: 			## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: cp-local-bin
cp-local-bin:
	mkdir -p $(DESTDIR)/$(prefix)/bin 
	cp -rf bin/* $(DESTDIR)/$(prefix)/bin
	# cp -rf scripts $(DESTDIR)/$(prefix)/scripts
	# cp -rf ./tmp/static $(DESTDIR)/$(prefix)/
