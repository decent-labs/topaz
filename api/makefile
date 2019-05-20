GOCMD=go
BUILDDIR=build
BINNAME=topaz-api

GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

build:
	$(GOBUILD) -mod=vendor -o ../$(BUILDDIR)/$(BINNAME) -v

clean:
	$(GOCLEAN)
	rm -rf ../$(BUILDDIR)/$(BINNAME)*

run: build
	cd .. && ./$(BUILDDIR)/$(BINNAME)
