<!-- QED cleanups TODO Needs to be re-thought and moved after structure
   (Being) what it was required to do. => Q. E. F.
   (Being) what it was required to prove. => Q. E. D.
   ...
-->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

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

  <!-- Removals -->
  <xsl:template match="div4[@type='Proof']/p[text() = 'Therefore etc. ']"/>
  <xsl:template match="div4[@type='Proof']/p[text() = 'Therefore etc.']"/>
  <!-- Prop I.8 -->
  <xsl:template match="div4[@type='Proof']/p[text() = 'If therefore etc.']"/>
  <!-- Prop I.15 -->
  <xsl:template match="div4[@type='Proof']/p[text() = 'Therefore etc. Q. E. D. ']"/>
  <!-- Prop I.36 -->
  <xsl:template match="div4[@type='Proof']/p[text() = 'Therefore etc. Q. E. D.']"/>

  <!-- Prop I.42 has a trailing Q. E. F. Apparently text() returns the array of
      interspersed text nodes and supports negative indexing. -->
  <xsl:template match="div3[@id='elem.1.42']/div4[@type='Proof']/p[position() = last()]">
    <p><xsl:value-of select="replace(text()[-1], ' Q. E. F.', '')"/></p>
  </xsl:template>

</xsl:stylesheet>
