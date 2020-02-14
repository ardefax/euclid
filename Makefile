V1 := data/heath/book01.json data/heath/book02.json
V2 := data/heath/book03.json data/heath/book04.json data/heath/book05.json data/heath/book06.json data/heath/book07.json data/heath/book08.json data/heath/book09.json
V3 := data/heath/book10.json data/heath/book11.json data/heath/book12.json data/heath/book13.json
JSON := $(V1) $(V2) $(V3)
MD := $(JSON:data/heath/book%.json=content/book%.md)

all: $(MD)

# Stubs out the content files with the proper data name
content/book%.md: data/heath/book%.json
	echo "---" > $@
	echo "data: $(@:content/%.md=%)" >> $@
	echo "type: book" >> $@
	echo "---" >> $@

# Intermediate trick from here https://stackoverflow.com/a/10609434
$(JSON): intermediate.json ;
.INTERMEDIATE: intermediate.json

intermediate.json: heath/books.xml book/book
	./book/book -d data/heath $<

book/book: book/main.go book/dom.go
	go build -o $@ ./book

.PHONY: clean
clean:
	rm -rf $(MD) $(JSON) heath/x/*.xml book/book


### XSLT Transforms

# brew install saxon
# http://www.saxonica.com/documentation/index.html#!using-xsl/commandline

N:=6

# Retain the intermediate files for inspection
.SECONDARY:

# Re-runs the identity transform
heath/books.xml: heath/x/$(N).xml heath/x/0.xslt
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
heath/x/0.xml: dl/books.xml heath/x/0.xslt
	saxon -s:$< -xsl:heath/x/0.xslt -o:$@
