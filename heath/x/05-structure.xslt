<!-- Adding the same structure to Books 2+ that exist in Book 1 -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="div1[@n != '1']/div2[starts-with(@n, 'Prop')]/div3">
    <xsl:copy>
      <xsl:apply-templates select="@*"/>
      <xsl:copy-of select="./head"/>
      <div4 type="Enunc">
        <xsl:copy-of select="./p[position() = 1]" />
      </div4>
      <div4 type="Proof">
        <xsl:copy-of select="./p[position() != 1 and position() != last()]" />
      </div4>
      <div4 type="QED">
        <xsl:copy-of select="./p[position() = last()]" />
      </div4>
    </xsl:copy>
  </xsl:template>

  <!-- Prop 1.40 was also poorly structured -->
  <xsl:template match="div3[@id='elem.1.40']">
    <xsl:copy>
      <xsl:apply-templates select="@*"/>
      <head>Proposition 40.</head>
      <div4 type="Enunc">
        <p><xsl:copy-of select="./p[position() = 1]/emph/text()" /></p>
      </div4>
      <div4 type="Proof">
        <xsl:copy-of select="./p[position() != 1 and position() != last()]" />
      </div4>
      <div4 type="QED">
        <p>Q. E. D.</p>
      </div4>
    </xsl:copy>
  </xsl:template>

</xsl:stylesheet>
