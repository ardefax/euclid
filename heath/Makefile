# brew install saxon
# http://www.saxonica.com/documentation/index.html#!using-xsl/commandline

N:=6

# Retain the intermediate files for inspection
.SECONDARY:

.PHONY: all
all: books.xml

# Re-runs the identity transform
books.xml: x/$(N).xml x/0.xslt
	saxon -s:$< -xsl:x/0.xslt -o:$@

x/6.xml: x/5.xml x/6.xslt
	saxon -s:$< -xsl:x/6.xslt -o:$@
x/5.xml: x/4.xml x/5.xslt
	saxon -s:$< -xsl:x/5.xslt -o:$@
x/4.xml: x/3.xml x/4.xslt
	saxon -s:$< -xsl:x/4.xslt -o:$@
x/3.xml: x/2.xml x/3.xslt
	saxon -s:$< -xsl:x/3.xslt -o:$@
x/2.xml: x/1.xml x/2.xslt
	saxon -s:$< -xsl:x/2.xslt -o:$@
x/1.xml: x/0.xml x/1.xslt
	saxon -s:$< -xsl:x/1.xslt -o:$@
x/0.xml: dl/books.xml x/0.xslt
	saxon -s:$< -xsl:x/0.xslt -o:$@

.PHONY: clean
clean:
	rm -rf x/*.xml
