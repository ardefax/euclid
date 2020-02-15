# Makefile that drives building of the site content.
#
# 	heath/dl/books.xml => ... => content/books/*.md
#
# Note that some of these build artifacts are checked in
# easily diff changes over time.

# Targets for each book
BOOKS := content/books/1.md \
         content/books/2.md \
         content/books/3.md \
         content/books/4.md \
         content/books/5.md \
         content/books/6.md \
         content/books/7.md \
         content/books/8.md \
         content/books/9.md \
         content/books/10.md \
         content/books/11.md \
         content/books/12.md \
         content/books/13.md

.PHONY: all 
all: $(BOOKS)

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

.PHONY: clean
clean:
	rm -rf $(BOOKS) heath/books.xml heath/x/*.xml book/book

