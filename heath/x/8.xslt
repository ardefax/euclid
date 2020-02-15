<!-- Building new structure and fixing other errors -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- TODO Figure out what I want to do with the headings, e.g. periods and caps -->
  <!-- Stripping Book # prefixes from propositions -->
  <xsl:template match="div2[@n = 'Prop']/head[starts-with(text(), 'BOOK')]">
    <head>PROPOSITIONS.</head>
  </xsl:template>

</xsl:stylesheet>
