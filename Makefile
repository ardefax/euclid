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
content: heath/books.xml book/book
	./book/book -d content/books $<


# Build the go binary used to generate content
book/book: book/main.go book/dom.go
	go build -o $@ ./book


### XSLT Transforms of the source material

# Re-runs the identity transform
.SECONDARY: heath/books.xml
heath/books.xml: heath/x/8.xml heath/x/0.xslt
	saxon -s:$< -xsl:heath/x/0.xslt -o:$@
heath/x/8.xml: heath/x/7.xml heath/x/8.xslt
	saxon -s:$< -xsl:heath/x/8.xslt -o:$@
heath/x/7.xml: heath/x/6.xml heath/x/7.xslt
	saxon -s:$< -xsl:heath/x/7.xslt -o:$@
heath/x/6.xml: heath/x/5.xml heath/x/6.xslt
	saxon -s:$< -xsl:heath/x/6.xslt -o:$@
heath/x/5.xml: heath/x/4.xml heath/x/5.xslt
	saxon -s:$< -xsl:heath/x/5.xslt -o:$@
heath/x/4.xml: heath/x/3.xml heath/x/4.xslt
	saxon -s:$< -xsl:heath/x/4.xslt -o:$@
heath/x/3.xml: heath/x/2.xml heath/x/3.xslt
	saxon -s:$< -xsl:heath/x/3.xslt -o:$@
heath/x/2.xml: heath/x/1.xml heath/x/2.xslt
	saxon -s:$< -xsl:heath/x/2.xslt -o:$@
heath/x/1.xml: heath/x/0.xml heath/x/1.xslt
	saxon -s:$< -xsl:heath/x/1.xslt -o:$@
heath/x/0.xml: heath/dl/books.xml heath/x/0.xslt
	saxon -s:$< -xsl:heath/x/0.xslt -o:$@

.PHONY: clean
clean:
	rm -rf $(BOOKS) heath/books.xml heath/x/*.xml book/book

