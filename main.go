// build extracts content from the Heath version
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
			fmt.Printf("%#v\n", d)
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

type Div1 struct {
	N      string `xml:"n,attr"`
	Type   string `xml:"type,attr"`
	Org    string `xml:"org,attr"`
	Sample string `xml:"sample,attr"`
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

