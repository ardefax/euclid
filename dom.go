package main

import (
	"encoding/xml"
	"fmt"
)

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
