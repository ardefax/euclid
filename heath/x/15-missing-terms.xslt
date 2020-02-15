<!-- Wrap definition words with the <term> tag where missing
  TODO There's probably a function I could define for this
-->
<xsl:stylesheet version="3.0"
  xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:mode on-no-match="shallow-copy" />

  <!-- https://www.w3.org/TR/xslt-30/#regex-examples -->
  <!-- Book 1 -->
  <xsl:template match="div3[@id='elem.1.def.19']/p/text()">
    <xsl:analyze-string select="." regex="Rectilineal|trilateral|quadrilateral|multilateral">
      <xsl:matching-substring>
        <term><xsl:value-of select="."/></term>
      </xsl:matching-substring>
      <xsl:non-matching-substring>
        <xsl:value-of select="."/>
      </xsl:non-matching-substring>
    </xsl:analyze-string>
  </xsl:template>

<!--
§Definition 20
Of trilateral figures, an equilateral triangle is that which has its three sides equal, an isosceles triangle that which has two of its sides alone equal, and a scalene triangle that which has its three sides unequal.

§Definition 21
Further, of trilateral figures, a right-angled triangle is that which has a right angle, an obtuse-angled triangle that which has an obtuse angle, and an acuteangled triangle that which has its three angles acute.

§Definition 22
Of quadrilateral figures, a square is that which is both equilateral and right-angled; an oblong that which is right-angled but not equilateral; a rhombus that which is equilateral but not right-angled; and a rhomboid that which has its opposite sides and angles equal to one another but is neither equilateral nor right-angled. And let quadrilaterals other than these be called trapezia.

§Definition 23
Parallel straight lines are straight lines which, being in the same plane and being produced indefinitely in both directions, do not meet one another in either direction.
-->

</xsl:stylesheet>
