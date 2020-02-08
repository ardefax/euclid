// build extracts content from the Heath version
//
// TODO:
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

		for _, d1 := range vol.Text.Body.Divs {
			debug("%#v\n", d1.div)
			if d1.Type != "book" {
				log.Fatal("invalid d1.type: %q (book)", d1.Type)
			}

			var bk Book
			for _, d2 := range d1.Divs {
				fmt.Printf("  %#v\n", d2.div)
				if d2.Type != "type" {
					log.Fatal("invalid d2.type: %q (book)", d1.Type)
				}

				switch d2.Type {
				case "type":
					switch d2.N {
						// TODO Book X interleaves definitions and propostions
						// in sections as `Def #` or `Prop #`
						case "Def":
							bk.Definitions = defs(d2)
						case "Post":
							bk.Postulates = posts(d2)
						case "CN":
							bk.CommonNotions = cns(d2)
						case "Prop":
							bk.Propositions = props(d2)
						default:
							log.Fatalf("invalid d2.N: %q (Def|Post|CN|Prop)", d2.N)
					}
				default:
					log.Fatalf("invalid type: %q (type)", d2.Type)
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
// raw exposes unprocessed XML data
type raw struct {
    Content []byte     `xml:",innerxml"`
    Nodes   []Node     `xml:",any"`
}
type Div1 struct {
	Divs []Div2 `xml:"div2"`
	div
}
type Div2 struct {
	Divs []Div3 `xml:"div3"`
	div
}
type Div3 struct {
	ID string `xml:"id,attr"`
	Divs []Div4 `xml:"div4"`
	div
	raw
}
// Looks like only book 1 uses the div4 structure to differentate between
// the Proposition statement (type="Enunc", e.g. Enunciation) and the Proof steps (type="Proof")
// although some of the later books propositions that nest div4 with (type="porism")
// and (type="lemma")
type Div4 struct {
	div
	raw
}

// TODO Later books (e.g. X) interleave definitions and propositions
type Book struct {
	Definitions []Definition
	Postulates []Postulate
	CommonNotions []CommonNotion
	Propositions []Proposition
}

type Definition string
func defs(d2 Div2) []Definition {
	a := make([]Definition, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:definition)", d3.Type)
		}
		a[i] = Definition(d3.Content)
		fmt.Println("d3:%s %s", d3.ID, a[i])
	}
	return a
}

type Postulate string
func posts(d2 Div2) []Postulate {
	a := make([]Postulate, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:postulate)", d3.Type)
		}
		a[i] = Postulate(d3.Content)
		debug("d3:%s: %s", d3.ID, a[i])
	}
	return a
}

type CommonNotion string
func cns(d2 Div2) []CommonNotion {
	a := make([]CommonNotion, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:common-notion)", d3.Type)
		}
		a[i] = CommonNotion(d3.Content)
		debug("d3:%s: %s", d3.ID, a[i])
	}
	return a
}

type Proposition struct {
	Claim string // TODO Enunciation?
	// TODO Construction, etc.
	Proof string // TODO Steps

	// TODO For tne non-div4 cases can we use the p tags?
	Raw string
}
func props(d2 Div2) []Proposition {
	a := make([]Proposition, len(d2.Divs))
	for i, d3 := range d2.Divs {
		var prop Proposition

		// TODO Book II also uses type="proposition"
		if d3.Type != "number" && d3.Type != "proposition" {
			log.Fatalf("invalid d3.type: %q (number:proposition)", d3.Type)
		}

		for _, d4 := range d3.Divs {
			switch d4.Type {
			case "Enunc":
				prop.Claim = string(d4.Content)
			case "Proof":
				prop.Proof = string(d4.Content)
			case "QED": // skip
			case "porism": // TODO
			default:
				log.Fatalf("invalid d4.type: %q (Enunc|Proof|QED|porism)", d4.Type)
			}
		}

		if prop.Proof == "" {
			// TODO Can we assume d3.Nodes[0] as Claim and d3.Nodes[1:] as Proof?
			prop.Raw = string(d3.Content)
			debug("d3:%s:raw %s", d3.ID, prop.Raw)
		} else {
			debug("d3:%s:claim %s\n", d3.ID, prop.Claim)
			debug("d3:%s:proof %s\n", d3.ID, prop.Proof)
		}
		a[i] = prop
	}
	return a
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
