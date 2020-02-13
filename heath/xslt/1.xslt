<!-- Strip out line break (lb), page break (pg), and figure tags -->
<!-- TODO Ideally should clean the stray whitespace from pb removals -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="lb"/>
  <xsl:template match="pb"/>
  <xsl:template match="figure"/>

</xsl:stylesheet>
