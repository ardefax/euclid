// build extracts content from the Heath version
//
// TODO:
//	Mapping onto better data types?
//	Rewriting this is a yaml or JSON doc...
//	Differentiating Propositions from Proof steps beyond book 1
package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	verbose = flag.Bool("v", false, "verbose debug spew")
	dir = flag.String("d", "data", "output directory")
)

func main() {
	log.SetPrefix(os.Args[0] + ":")
	log.SetFlags(0)
	flag.Parse()
	if flag.NArg() == 0 || *dir == "" {
		fmt.Fprintf(os.Stderr, "usage: %s [options] <volumes>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	for _, arg := range flag.Args() {
		debug("reading: %s", arg)
		buf, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}

		var vol Volume
		err = xml.Unmarshal(buf, &vol)
		if err != nil {
			log.Fatal(err)
		}
		for _, d := range vol.Text.Front.Divs {
			fmt.Printf("%#v\n", d)
		}
		for _, d := range vol.Text.Body.Divs {
			fmt.Printf("%#v\n", d.div)
			for _, d2 := range d.Divs {
				fmt.Printf("  %#v\n", d2.div)
				for _, d3 := range d2.Divs {
					fmt.Printf("    %s %#v\n", d3.ID, d3.div)
					for _, d4 := range d3.Divs {
						fmt.Printf("      %#v\n", d4.div)
					}
				}
			}
		}
	}
}

type Volume struct {
	XMLName xml.Name `xml:"TEI.2"`
	Text Text `xml:"text"`
}

type Text struct {
	Front Front `xml:"front"`
	Body Body `xml:"body"`
}

type Front struct {
	Divs []Div1 `xml:"div1"`
}

type Body struct {
	Divs []Div1 `xml:"div1"`
}

// div are common values across Div# elements.
type div struct {
	N      string `xml:"n,attr"`
	Type   string `xml:"type,attr"`
	Org    string `xml:"org,attr"`
	Sample string `xml:"sample,attr"`

	Heads []string `xml:"head"`
}

type Div1 struct {
	div
	Divs []Div2 `xml:"div2"`
}
type Div2 struct {
	div
	Divs []Div3 `xml:"div3"`
}
type Div3 struct {
	ID string `xml:"id,attr"`
	div
	Divs []Div4 `xml:"div4"`
}
// Looks like only book 1 uses the div4 structure to differentate between
// the Proposition statement (type="Enunc") and the Proof steps (type="Proof")
// although some of the later books propositions that nest div4 with (type="porism")
// and (type="lemma")
type Div4 struct {
	div
}

func debug(format string, head interface{}, tail ...interface{}) {
	if !*verbose {
		return
	}
	format = os.Args[0] + ": " + format + "\n"
	args := make([]interface{}, 1, len(tail) + 1)
	args[0] = head
	args = append(args, tail...)
	fmt.Fprintf(os.Stderr, format, args...)
}

