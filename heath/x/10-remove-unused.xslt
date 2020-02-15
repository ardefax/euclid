<!-- Strip out line break (lb), page break (pg), figure, and some hi tags -->
<!-- TODO Ideally should clean the stray whitespace from pb removals -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <xsl:template match="lb"/>
  <xsl:template match="pb"/>
  <xsl:template match="figure"/>

  <!-- TODO these render tags are interesting -->
  <xsl:template match="hi"> <!-- TODO: Value in keeping around? -->
      <xsl:copy-of select="*|text()"/>
  </xsl:template>

</xsl:stylesheet>
