<!-- QED cleanups
   (Being) what it was required to do. => Q. E. F.
   (Being) what it was required to prove. => Q. E. D.
-->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!--
  <xsl:template match="div4[contains(text(), 'required to do')]">
    <xsl:copy>
      <xsl:copy-of select="@*"/>
      <xsl:text>Q.E.F.</xsl:text>
    </xsl:copy>
  </xsl:template>
  -->

  <xsl:template match="div4[@type='QED']/p[text() = '(Being) what it was required to prove.']">
    <xsl:copy>
      <xsl:copy-of select="@*"/>
      <xsl:text>Q. E. D.</xsl:text>
    </xsl:copy>
  </xsl:template>
  <xsl:template match="div4[@type='QED']/p[text() = '(Being) what it was required to do.']">
    <xsl:copy>
      <xsl:copy-of select="@*"/>
      <xsl:text>Q. E. F.</xsl:text>
    </xsl:copy>
  </xsl:template>
</xsl:stylesheet>
