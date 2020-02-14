<!-- Replace term => dfn and emph => var in the body -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="body//term">
    <dfn><xsl:apply-templates select="@* | node()"/></dfn>
  </xsl:template>
  <xsl:template match="body//emph">
    <var><xsl:apply-templates select="@* | node()"/></var>
  </xsl:template>
</xsl:stylesheet>
