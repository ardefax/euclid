<!-- Building new structure and fixing other errors -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <!-- TODO prints attributes on there own lines -->
  <!-- <xsl:output method="xml" indent="yes" /> -->
  <xsl:mode on-no-match="shallow-copy" />

  <!-- TODO Prop I.37.7 and I.38.7 are bracketed with a note (e.g should reference CN.3)
  <xsl:template match="div3[@id='elem.1.37' or id='elem.1.38']/div4[@type=Proof]/p[7]">
  </xsl:template>
  -->

  <!-- Prop 1.40 was poorly structured
      https://stackoverflow.com/a/5672132 -->
  <xsl:template match="div3[@id='elem.1.40']">
    <xsl:copy>
      <xsl:apply-templates select="@*"/>
      <head>Proposition 40.</head>
      <div4 type="Enunc">
        <p><xsl:copy-of select="./p[position() = 1]/var/text()" /></p>
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
