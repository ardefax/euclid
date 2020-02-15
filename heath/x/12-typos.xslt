<!-- Typo fixes that can be resolved by simple regex -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  xmlns:ardefax="http://xsl.ardefax.org">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="div3[@id='elem.1.def.21']/p/text()">
    <xsl:analyze-string select="." regex="acuteangled">
      <xsl:matching-substring>acute-angled</xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

</xsl:stylesheet>
