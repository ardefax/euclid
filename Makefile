# Makefile that drives building of the site content.
#
# 	heath/dl/books.xml => ... => content/books/**
#
# Note that some of these build artifacts are checked in
# easily diff changes over time.

# Targets for each book
BOOKS := \
  content/books/1/_index.md \
  content/books/2/_index.md \
  content/books/3/_index.md \
  content/books/4/_index.md \
  content/books/5/_index.md \
  content/books/6/_index.md \
  content/books/7/_index.md \
  content/books/8/_index.md \
  content/books/9/_index.md \
  content/books/10/_index.md \
  content/books/11/_index.md \
  content/books/12/_index.md \
  content/books/13/_index.md

.PHONY: all overview
all: $(BOOKS) overview

overview: content/books/_index.md
content/books/_index.md:
	@echo --- > $@
	@echo title: Overview >> $@
	@echo --- >> $@

# Intermediate trick for multiple targets from single recipe
# https://stackoverflow.com/a/10609434
$(BOOKS): content ;
.INTERMEDIATE: content
# TODO This runs even if everything is up to date...
content: heath/books.xml book/book
	./book/book -d content/books $<


# Build the go binary used to generate content
book/book: book/main.go book/dom.go
	go build -o $@ ./book

### XSLT Transforms of the source material

# TODO Figure out if it's feasible to create an implicit
# rule that relies on the sequential nature of the files.
# Alternatively, could use an automatic variable trick to
# figure out the "lowest-valued" xslt that has changed
# as the start of what needs to be run.
TRANSFORMS := $(wildcard heath/x/*.xslt)
heath/books.xml: heath/xslt.sh heath/dl/books.xml $(TRANSFORMS)
	./heath/xslt.sh heath/dl/books.xml heath/x heath/books.xml

# TODO: Run these through an actual expansion since the inputs
# are minimal and currently require an initial JS step to figure
# out things like intersections, label positions, etc.
#content/books/1/prop.1-fig.svg: figures/b1.prop1.svg
#	mkdir -p $(@D); cp  $< $@

.PHONY: clean
clean:
	rm -rf content/books heath/books.xml tmp/* book/book

