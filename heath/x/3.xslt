<!-- A little more involved updating of refs to links 
  https://stackoverflow.com/questions/10184694#comment20903136_14902493 -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="body//ref">
    <a href="#{@target}"><xsl:value-of select="text()"/></a>
  </xsl:template>
</xsl:stylesheet>
