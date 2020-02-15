<!-- Wrap definition words with the <term> tag where missing -->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
  xmlns:ardefax="http://xsl.ardefax.org">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- https://www.w3.org/TR/xslt-30/#regex-examples -->
  <!-- Book 1 -->
  <xsl:template match="div3[@id='elem.1.def.19']/p/text()">
    <xsl:analyze-string select="." regex="Rectilineal|trilateral|quadrilateral|multilateral">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <xsl:template match="div3[@id='elem.1.def.20']/p/text()">
    <xsl:analyze-string select="." regex="(equilateral|isosceles|scalene) triangle">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <xsl:template match="div3[@id='elem.1.def.21']/p/text()">
    <xsl:analyze-string select="." regex="(right|obtuse|acute)-angled triangle">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <xsl:template match="div3[@id='elem.1.def.22']/p/text()">
    <xsl:analyze-string select="." regex="(square|oblong|rhombus|rhomboid|trapezia)">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <xsl:template match="div3[@id='elem.1.def.23']/p/text()">
    <xsl:analyze-string select="." regex="(Parallel straight lines)">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <!-- Book 2 -->
  <xsl:template match="div3[@id='elem.2.def.1']/p/text()">
    <xsl:analyze-string select="." regex="(contained)">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <xsl:template match="div3[@id='elem.2.def.2']/p/text()">
    <xsl:analyze-string select="." regex="(gnomon)">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring><xsl:value-of select="."/></xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

  <!-- Book 3 -->




  <!-- TODO Can't quite seem to make this work as a function (maybe a for-each eventually)
  <xsl:function name="ardefax:find-terms">
    <xsl:param name="doc"/>
    <xsl:param name="id"/>
    <xsl:param name="terms"/>
    <xsl:analyze-string select="$doc//div3[@id=$id]/p/text()" regex="$terms">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring>
        <xsl:value-of select="."/>
      </xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:function>

  <xsl:template match="/">
    <xsl:value-of select="ardefax:find-terms(., 'elem.1.def.19', 'Rectilineal|trilateral|quadrilateral|multilateral')"/>
  </xsl:template>
  -->

</xsl:stylesheet>
