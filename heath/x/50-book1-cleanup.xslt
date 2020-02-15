<!-- Building new structure and fixing other errors -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- I.Post.1 remove the first line that will be "promoted" in 6.xslt -->
  <xsl:template match="div3[@id='elem.1.post.1']/p[position() = 1]" />

  <!-- I.CN.4 and I.CN.5 - Remove the stray footnotes with unclear references -->
  <xsl:template match="div3[@id='elem.1.c.n.4' or @id='elem.1.c.n.5']/p">
    <p><xsl:value-of select="substring(text(), 5)"/></p>
  </xsl:template>

  <!-- Prop I.15, I.36 missing the QED div. The copy-of `.` causes issues when
      this was also attempted to be combined with the prior QED fixes -->
  <xsl:template match="div3[@id='elem.1.15' or @id='elem.1.36']/div4[@type='Proof']">
    <xsl:copy-of select="."/>
    <div4 type="QED" org="uniform" sample="complete">
      <p>Q. E. D.</p>
    </div4>
  </xsl:template>
  <!-- Prop I.42 is missing QED div -->
  <xsl:template match="div3[@id='elem.1.42']/div4[@type='Proof']">
    <xsl:copy-of select="."/>
    <div4 type="QED" org="uniform" sample="complete">
      <p>Q. E. F.</p>
    </div4>
  </xsl:template>

  <!-- Prop I.34 random missing space. -->
  <xsl:template match="div3[@id='elem.1.34']/div4[@type='Proof']/p[position() = 3]/var[position() = 3]">
    <var>BC</var><xsl:text> </xsl:text>
  </xsl:template>

</xsl:stylesheet>
