#/usr/bin/env bash
#
# usage:
#   xslt.sh <in>.xml <dir> <out>.xml
#
# Sequentially process an input xml file with a directory of
# xslt files to produce a new xml file. The intermediate files
# are left in a `tmp` directory with the same basename of the
# source xslt files.
#

set -eu

in="${1}"
dir="${2}"
out="${3}"

for xslt in ${dir}/*.xslt
do
    xml="tmp/$(basename "${xslt}" .xslt).xml"
	echo "saxon -s:${in} -xsl:${xslt} -o:${xml}"
	saxon -s:${in} -xsl:${xslt} -o:${xml}
    in=${xml}
done

echo "mv ${in} ${out}"
mv ${in} ${out}
