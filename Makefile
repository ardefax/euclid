# Targets for each book
MD := content/books/1.md \
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
all: $(MD)

# Intermediate trick from here https://stackoverflow.com/a/10609434
$(MD): intermediate.md ;
.INTERMEDIATE: intermediate.md
intermediate.md: heath/books.xml book/book
	./book/book -d content/books $<


# Build the book binary that drives content gen
book/book: book/main.go book/dom.go
	go build -o $@ ./book


### XSLT Transforms of the source material

# Re-runs the identity transform
.SECONDARY: heath/books.xml
heath/books.xml: heath/x/6.xml heath/x/0.xslt
	saxon -s:$< -xsl:heath/x/0.xslt -o:$@
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
	rm -rf $(MD) $(JSON) heath/x/*.xml book/book

