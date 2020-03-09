# Heath

Contains the XML text for Heath's version of Euclid's Elements and tool(s) to
transform it. The original files are in `dl/vol{1,2,3}.xml` along with the merged
`dl/books.xml` which excludes the initial chapters.

The `x` directory exists to do various transforms on the source text that _should_
work on both `books.xml` and `vol?.xml`, however, only the `books.xml` is consumed.

The above transforms are driven by the `Makefile` in the root of the repo using
the [saxon] XSLT transformer via the `xslt.sh` script in this directory.

[saxon]: http://www.saxonica.com/documentation/index.html#!using-xsl/commandline
