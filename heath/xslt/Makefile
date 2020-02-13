# brew install saxon
# http://www.saxonica.com/documentation/index.html#!using-xsl/commandline

.PHONY: all
all: 3/vol1.xml 3/vol2.xml 3/vol3.xml

3/vol1.xml: 2/vol1.xml 3.xslt 
	saxon -s:$< -xsl:3.xslt -o:$@
3/vol2.xml: 2/vol2.xml 3.xslt 
	saxon -s:$< -xsl:3.xslt -o:$@
3/vol3.xml: 2/vol3.xml 3.xslt 
	saxon -s:$< -xsl:3.xslt -o:$@

2/vol1.xml: 1/vol1.xml 2.xslt 
	saxon -s:$< -xsl:2.xslt -o:$@
2/vol2.xml: 1/vol2.xml 2.xslt 
	saxon -s:$< -xsl:2.xslt -o:$@
2/vol3.xml: 1/vol3.xml 2.xslt 
	saxon -s:$< -xsl:2.xslt -o:$@

1/vol1.xml: 0/vol1.xml 1.xslt 
	saxon -s:$< -xsl:1.xslt -o:$@
1/vol2.xml: 0/vol2.xml 1.xslt 
	saxon -s:$< -xsl:1.xslt -o:$@
1/vol3.xml: 0/vol3.xml 1.xslt 
	saxon -s:$< -xsl:1.xslt -o:$@

0/vol1.xml: ../vol1.xml 0.xslt 
	saxon -s:$< -xsl:0.xslt -o:$@
0/vol2.xml: ../vol2.xml 0.xslt 
	saxon -s:$< -xsl:0.xslt -o:$@
0/vol3.xml: ../vol3.xml 0.xslt 
	saxon -s:$< -xsl:0.xslt -o:$@

.PHONY: clean
clean:
	rm -rf 0/* 1/* 2/* 3/*
