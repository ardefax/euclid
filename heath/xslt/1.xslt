<!-- Strip out line break (lb), page break (pg), figure, and hi tags -->
<!-- TODO Ideally should clean the stray whitespace from pb removals -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="lb"/>
  <xsl:template match="pb"/>
  <xsl:template match="figure"/>

  <!-- Want the embedded text and child nodes to be retained -->
  <xsl:template match="hi">
      <xsl:copy-of select="*|text()"/>
  </xsl:template>


</xsl:stylesheet>
