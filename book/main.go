// book extracts content from the Heath version
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
		fmt.Fprintf(os.Stderr, "usage: %s [options] <books.xml>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	books := make([]*Book, 0)
	for _, arg := range flag.Args() {
		debug("reading: %s", arg)
		buf, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}

		// Decode the raw XML structure
		var body Body
		err = xml.Unmarshal(buf, &body)
		if err != nil {
			log.Fatal(err)
		}

		// Transfrom to the more strucuted Book type.
		for _, d1 := range body.Divs {
			book := new(Book)
			if err := book.parse(d1); err != nil {
				log.Fatal(err)
			}
			books = append(books, book)
		}
	}

	if err := os.MkdirAll(*dir, 0755); err != nil {
		log.Fatal(err)
	}
	for _, b := range books {
		// Write the JSON as the markdown file frontmatter
		// The layouts will do the heavy lifting of content generation
		path := fmt.Sprintf(filepath.Join(*dir, "%d.md"), b.Number)
		f, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}

		enc := json.NewEncoder(f)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		debug("writing %s", path)
		if err := enc.Encode(b); err != nil {
			log.Fatal(err)
		}
	}
}

// Book is one of Euclid's Elements books.
type Book struct {
	// Title of the book
	Title string `json:"title"`
	// Number of the book
	Number int `json:"number"`
	// Roman numeral string of the book number
	Roman string `json:"roman"`
	// Sections from the content and structure of a book
	Sections []Section `json:"sections"`

	// Hugo stuff

	// Weight is the book Number, used by Hugo for sorting.
	Weight int `json:"weight"`
}

// Section is a generic part of the book
type Section struct {
	// ID is used to uniquely referenece a section. Can be suffixed
	// with an index to reference a specific text paragraph.
	ID string `json:"id"`
	// Kind is is the type of section
	Kind string `json:"kind"`
	// Title is used for  section headings
	Title string `json:"title"`
	// Text is a list of paragraphs that may contain embedded HTML
	Text []string `json:"text"`
	// Sections are child sections, rendered after the above text.
	Sections []Section `json:"sections"`
}

func (b *Book) parse(d1 Div1) error {
	debug("book: %#v\n", d1.div)
	if d1.Type != "book" {
		return fmt.Errorf("invalid d1.type: %q (book)", d1.Type)
	}

	n, err := strconv.Atoi(d1.N)
	if err != nil {
		return fmt.Errorf("invalid d1.N: %s", err)
	}
	b.Number = n
	b.Weight = n
	b.Roman = roman(n)
	b.Title = fmt.Sprintf("BOOK %s", b.Roman)

	for _, d2 := range d1.Divs {
		b.parseSection(d2)
	}
	return nil
}

// parseSection converts a Div2 to the equivalent section of the book
func (b *Book) parseSection(d2 Div2) error {
	debug("  %#v\n", d2.div)
	if d2.Type != "type" {
		return fmt.Errorf("invalid d2.type: %q (book)", d2.Type)
	}
	if len(d2.Heads) != 1 {
		warn("b%d.%s: wrong # of d2.heads: %d (1)", b.Number, d2.N, len(d2.Heads))
	}

	var s Section

	// Book X interleaves sections as `Def #` or `Prop #`
	short := strings.Replace(strings.ToLower(d2.N), " ", "", -1)
	switch short {
	case "def", "def1", "def2", "def3":
		s = b.parseSimple(d2, short, "definition")
	case "post":
		s = b.parseSimple(d2, short, "postulate")
	case "cn": // TODO the c.n format is lame
		s = b.parseSimple(d2, "c.n", "common-notion")
	case "prop", "prop1", "prop2", "prop3":
		s = b.parseProps(d2, short)
	default:
		return fmt.Errorf("elem.%d: invalid d2.N: %q (Def|Post|CN|Prop)", b.Number, d2.N)
	}

	b.Sections = append(b.Sections, s)
	return nil
}

// parseSimple converts a non-proposition Div2 into a Section
func (b *Book) parseSimple(d2 Div2, short, kind string) Section {
	switch kind {
	case "definition", "postulate", "common-notion":
		// no-op
	case "proposition":
		log.Fatalf("section: not simple kind: %s", kind)
	default:
		log.Fatalf("section: invalid kind: %s", kind)
	}

	s := Section{
		ID: fmt.Sprintf("elem.%d.%s", b.Number, short),
		Kind: "list:" + kind,
		Title: cleanContent(d2.Heads[0]),
		Sections: make([]Section, len(d2.Divs)),
	}

	title := strings.Title(kind)
	for i, d3 := range d2.Divs {
		if d3.Type != "number" { // TODO fix the proposition one
			log.Fatalf("invalid d3.type: %q (number:%s)", d3.Type, kind)
		}
		// TODO V.Def.17 has two paragraphs (and a couple others)
		// TODO This should just join paras into text but Postulate 1 needs to be fixed
		text := []string{cleanContent(d3.Paras[0])}
		if len(d3.Paras) != 1 {
			warn("%s: wrong # of d3.paras: %d (1:%s)", d3.ID, len(d3.Paras), kind)
			text = append(text, cleanContent(d3.Paras[1]))
		}
		s.Sections[i] = Section{
			ID: d3.ID,
			Kind: kind,
			Title: fmt.Sprintf("%s %s", title, d3.N),
			Text: text,
		}
		debug("d3:%s %s", d3.ID, s.Sections[i].Text[0])
	}
	return s
}


// parseProps TODO doc
func (b *Book) parseProps(d2 Div2, short string) Section {
	s := Section{
		ID: fmt.Sprintf("elem.%d.prop", b.Number),
		Kind: "list:proposition",
		Title: cleanContent(d2.Heads[0]),
		Sections: make([]Section, len(d2.Divs)),
	}
	for i, d3 := range d2.Divs {
		// XXX Book II also uses type="proposition"
		if d3.Type != "number" && d3.Type != "proposition" {
			log.Fatalf("invalid d3.type: %q (number:proposition)", d3.Type)
		}
		ss := Section{
			ID: d3.ID,
			Kind: "proposition",
			Title: fmt.Sprintf("Proposition %s.", d3.N),
		}
		for _, d4 := range d3.Divs {
			switch d4.Type {
			case "Enunc":
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the theorem, not %d", len(d4.Paras))
				}
				ss.Sections = append(ss.Sections, Section {
					ID: d3.ID + ".theorem", // TODO Rationalize ID strat better
					Kind: "theorem",
					Text: []string{ cleanContent(d4.Paras[0]) },
				})
			case "Proof":
				if len(d4.Paras) < 2 {
					log.Fatalf("Expected some steps for the proof, not %d", len(d4.Paras))
				}
				sss := Section{
					ID: d3.ID + ".proof", // TODO Rationalize ID strat better
					Kind: "proof",
					Text: make([]string, len(d4.Paras)),
				}
				for j, p := range d4.Paras {
					sss.Text[j] = cleanContent(p)
				}
				ss.Sections = append(ss.Sections, sss)
			case "QED":
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the QED, not %d", len(d4.Paras))
				}
				ss.Sections = append(ss.Sections, Section{
					Kind: "qed",
					Text: []string{cleanContent(d4.Paras[0])},
				})
			case "porism": // TODO
			case "lemma": // TODO
			default:
				log.Fatalf("invalid d4.type: %q (Enunc|Proof|QED|porism|lemma)", d4.Type)
			}
		}

		s.Sections[i] = ss
	}
	return s
}

// cleanContent does any content transformation that may be necessary on a node
// in the source text. For now it just unpacks the inner text content as-is.
func cleanContent(p Node) string {
	return string(p.Content)
}

// roman is a terrible implementation of making roman numerals
func roman(n int) string {
	switch n {
	case 1: return "I"
	case 2: return "II"
	case 3: return "III"
	case 4: return "IV"
	case 5: return "V"
	case 6: return "VI"
	case 7: return "VII"
	case 8: return "VIII"
	case 9: return "IX"
	case 10: return "X"
	case 11: return "XI"
	case 12: return "XII"
	case 13: return "XIII"
	}
	log.Fatalf("TODO: proper roman numerals %d", n)
	return ""
}

func warn(format string, head interface{}, tail ...interface{}) {
	format = os.Args[0] + ": warn: " + format + "\n"
	args := make([]interface{}, 1, len(tail) + 1)
	args[0] = head
	args = append(args, tail...)
	fmt.Fprintf(os.Stderr, format, args...)
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
