package main

import (
	"encoding/xml"
	"fmt"
)

// Volume is the root element of the Heath text
type Volume struct {
	XMLName xml.Name `xml:"TEI.2"`
	Text Text `xml:"text"`
}
// Text wraps the front and body elements
type Text struct {
	Front Front `xml:"front"`
	Body Body `xml:"body"`
}
// Front is the element containting the commentary chapters.
type Front struct {
	Divs []Div1 `xml:"div1"`
}
// Body is the element containting books of the elements
type Body struct {
	Divs []Div1 `xml:"div1"`
}
// div are common values across Div# elements.
type div struct {
	N      string `xml:"n,attr"`
	Type   string `xml:"type,attr"`
	//TODO Org    string `xml:"org,attr"`
	//TODO Sample string `xml:"sample,attr"`

	//TODO Heads []string `xml:"head"`
}
// raw exposes unprocessed XML data
type raw struct {
    Content []byte     `xml:",innerxml"`
    Nodes   []Node     `xml:",any"`
}
// Div1 is a top-level book or chapter element.
type Div1 struct {
	Divs []Div2 `xml:"div2"`
	div
}
// Div2 is a nested element for grouping definitions, propositions, etc.
type Div2 struct {
	Divs []Div3 `xml:"div3"`
	div
}
// Div3 is a nested element for individual definitions, propositions, etc.
type Div3 struct {
	ID string `xml:"id,attr"`
	Divs []Div4 `xml:"div4"`
	div
	raw
}
// Div4 is a nested element used for sub-parts of propositions.
//
// Note: Only 'Book I' uses the div4 element to provide structure around the
// proposition statement (type="Enunc", e.g. Enunciation) and the Proof steps
// (type="Proof"). Later books have propositions with lemmas (type="lemma")
// and corollaries (type="porism").
type Div4 struct {
	div
	raw
}

// Node is a generic XML element.
type Node struct {
    XMLName xml.Name
    Attrs   []xml.Attr `xml:"-"`
    Content []byte     `xml:",innerxml"`
    Nodes   []Node     `xml:",any"`
}

// UnmarshalXML implements xml.Unmarshaler to also extract elements
func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    n.Attrs = start.Attr
    type node Node

    return d.DecodeElement((*node)(n), &start)
}

func (n Node) Dump() {
	n.Print("")
}

func (n Node) Print(prefix string) {
	fmt.Printf("%s%s\n", prefix, n.XMLName)
	for _, n := range n.Nodes {
		n.Print(prefix + " ")
	}
}
