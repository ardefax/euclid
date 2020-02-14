<!-- Building new structure and fixing other errors -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- Attempting a translation similar to Prop 1.40 for all Book 2
      https://stackoverflow.com/a/5672132 -->
  <xsl:template match="div1[@n='2']/div2[@n='Prop']/div3">
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

</xsl:stylesheet>
