# Go parameters
GOCMD=go
BUILDDIR=build
CURDIR=api
BINNAME=topaz-api

# GOTEST=$(GOCMD) test
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOGET=$(GOCMD) get

# all: test build

# test:
# 	$(GOTEST) -v ./...

build:
	$(GOBUILD) -o ../$(BUILDDIR)/$(BINNAME) -v

build-linux:
	$(GOGET) github.com/karalabe/xgo
	cd ../$(BUILDDIR) && xgo -out $(BINNAME) --targets=linux/amd64 ../$(CURDIR)/.

clean:
	$(GOCLEAN)
	rm -rf ../$(BUILDDIR)/$(BINNAME)*

run: build
	cd .. && ./$(BUILDDIR)/$(BINNAME)
