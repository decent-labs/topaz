GOCMD=go
BUILDDIR=build
BINNAME=topaz-batch

GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

build:
	$(GOBUILD) -o ../$(BUILDDIR)/$(BINNAME) -v

clean:
	$(GOCLEAN)
	rm -rf ../$(BUILDDIR)/$(BINNAME)*

run: build
	cd .. && ./$(BUILDDIR)/$(BINNAME)
