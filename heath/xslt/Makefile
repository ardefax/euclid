# brew install saxon
# http://www.saxonica.com/documentation/index.html#!using-xsl/commandline

N:=5

# Retain the intermediate files for inspection
.SECONDARY:

.PHONY: all
all: vol1.xml vol2.xml vol3.xml

vol%.xml: $(N)/vol%.xml
	ln -sf $<

5/vol%.xml: 4/vol%.xml 5.xslt
	saxon -s:$< -xsl:5.xslt -o:$@
4/vol%.xml: 3/vol%.xml 4.xslt
	saxon -s:$< -xsl:4.xslt -o:$@
3/vol%.xml: 2/vol%.xml 3.xslt
	saxon -s:$< -xsl:3.xslt -o:$@
2/vol%.xml: 1/vol%.xml 2.xslt
	saxon -s:$< -xsl:2.xslt -o:$@
1/vol%.xml: 0/vol%.xml 1.xslt
	saxon -s:$< -xsl:1.xslt -o:$@
0/vol%.xml: ../vol%.xml 0.xslt
	saxon -s:$< -xsl:0.xslt -o:$@

.PHONY: clean
clean:
	rm -rf `seq 0 $(N)` vol?.xml
