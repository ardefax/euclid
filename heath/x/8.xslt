<!-- Building new structure and fixing other errors -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- Post.1 First paragraph needs to be "promoted" -->
  <xsl:template match="div3[@n = '1']/div2[starts-with(@n, 'Prop')]/div3">
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
