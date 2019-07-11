GOCMD=go
BUILDDIR=build
BINNAME=topaz-batch
BINPATH=/$(BUILDDIR)/$(BINNAME)

GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

ifeq ($(OS),Windows_NT)
	BINPATH := $(subst /,\,$(BINPATH)).exe
	RMCMD = rd /s /q
else
	RMCMD = rm -rf
endif

build:
	$(GOBUILD) -o ..$(BINPATH) -v

clean:
	$(GOCLEAN)
	cd .. && $(RMCMD) $(BUILDDIR)

run: build
	cd .. && .$(BINPATH)
