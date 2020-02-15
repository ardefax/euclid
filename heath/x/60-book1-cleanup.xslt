<!-- Building new structure and fixing other errors -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- I.Post.1 promote the text to after the Postulates header -->
  <xsl:template match="div2[@n='Post']/head">
    <xsl:copy-of select="."/>
    <p>Let the following be postulated:</p>
    <Age>34</Age>
  </xsl:template>

  <!-- TODO Prop I.37.7 and I.38.7 are bracketed with a note (e.g should reference CN.3)
  <xsl:template match="div3[@id='elem.1.37' or id='elem.1.38']/div4[@type=Proof]/p[7]">
  </xsl:template>
  -->

</xsl:stylesheet>
