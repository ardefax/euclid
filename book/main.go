// book transforms the XML version of Heath's translation into the 
// content used to generate the site. The output structure has the
// following pattern
//
//	{{*dir}}/
//	  {{book.Number}}/
//	    _index.md
//		{{ range where book.Sections "Kind" "list:proposition" }}
//	      {{ range .Sections }} {{.frag}}.md {{ end }}
//      {{ end }}
//	
//
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
	"regexp"
	"strconv"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "verbose debug spew")
	dir = flag.String("d", "content", "output directory")
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

	setupNextPrevLinks(books)
	if err := writeBooksContent(books, *dir); err != nil {
		log.Fatal(err)
	}
}

func setupNextPrevLinks(books []*Book) {
	// Build out the prev/next links and book numbers TODO: improve?
	var prevBook *Book
	var prevSection *Section
	for _, b := range books {
		if prevSection != nil {
			// Last proposition points to next book
			prevSection.Next = &NavLink{
				Text: b.Title,
				URL: fmt.Sprintf("/books/%d", b.Number),
			}
		}
		prevSection = nil

		for _, s := range b.Sections {
			if s.Kind != "list:proposition" {
				continue
			}
			for _, ss := range s.Sections {
				ss.Book = b.Number
				if prevSection != nil {
					prevSection.Next = &NavLink {
						Text: ss.Title,
						URL: fmt.Sprintf("/books/%d/%s", ss.Book, ss.Frag),
					}
					ss.Prev = &NavLink{
						Text: prevSection.Title,
						URL: fmt.Sprintf("/books/%d/%s", ss.Book, prevSection.Frag),
					}
				} else {
					// First proposition points back to "this" book
					ss.Prev = &NavLink {
						Text: b.Title,
						URL: fmt.Sprintf("/books/%d", ss.Book),
					}
				}
				prevSection = ss
			}
		}

		if prevBook != nil {
			prevBook.Next = &NavLink{
				Text: b.Title,
				URL: fmt.Sprintf("/books/%d", b.Number),
			}
			b.Prev = &NavLink {
				Text: prevBook.Title,
				URL: fmt.Sprintf("/books/%d", prevBook.Number),
			}
		} else {
			// First book goes back to the overview
			b.Prev = &NavLink {
				Text: "Overview",
				URL: fmt.Sprintf("/books"),
			}
		}
		prevBook = b
	}
	// Link to the about page for the last section of the last book
	prevSection.Next = &NavLink{ Text: "About", URL: "/about", }
}

func writeBooksContent(books []*Book, dir string) error {
	for _, b := range books {
		root := filepath.Join(dir, strconv.Itoa(b.Number))
		if err := os.MkdirAll(root, 0755); err != nil {
			return err
		}
		for _, s := range b.Sections {
			if s.Kind != "list:proposition" {
				continue
			}
			for _, ss := range s.Sections {
				// TODO More hacking on what's written on prop vs book pages
				title := ss.Title
				ss.Title = fmt.Sprintf("BOOK %s: %s", b.Roman, title)
				if err := writeBookProp(ss, root); err != nil {
					return err
				}

				// TODO Kinda hacky so that we set title-links and don't emit
				// the non-theorem sections below
				ss.Title = title
				ss.Link = fmt.Sprintf("/books/%d/%s", b.Number, ss.Frag)
				keep := make([]*Section, 0)
				for _, sss := range ss.Sections {
					if sss.Kind == "theorem" {
						keep = append(keep, sss)
					}
				}
				ss.Sections = keep
			}
		}

		if err := writeBookIndex(b, root); err != nil {
			return err
		}
	}
	return nil
}

func writeBookIndex(b *Book, dir string) error {
	// Write the JSON as the markdown file frontmatter
	// The layouts will do the heavy lifting of content generation
	path := filepath.Join(dir, "_index.md")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	b.Type = "book" // TODO do this more consistently 
	b.Layout = "book" // TODO do this more consistently

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	debug("writing %s", path)
	return enc.Encode(b); // TODO Remove proof/QED text here instead of in CSS
}

func writeBookProp(s *Section, dir string) error {
	// Write the JSON as the markdown file frontmatter
	// The layouts will do the heavy lifting of content generation
	path := filepath.Join(dir, s.Frag + ".md")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	s.Layout = s.Kind // TODO do this more consistently 

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	debug("writing %s", path)
	return enc.Encode(s);
}

// NavLink is a wrapper for a text label and URL.
type NavLink struct {
	// Text is the displayed content of the link, e.g. inside the <a>{{.Text}}</a>
	Text string `json:"text"`
	// URL is the link reference, e.g. <a href="{{.URL}}">
	URL string `json:"url"`
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
	Sections []*Section `json:"sections"`

	// Hugo stuff

	// Weight is the book Number, used by Hugo for sorting.
	Weight int `json:"weight"`
	// Type is used to filter top-level books, set to "book"
	Type string `json:"type,omitempty"`
	// Layout is used to pick the right template file, set to "book"
	Layout string `json:"layout,omitempty"`
	// Next is for setting up navigation in hugo templates, e.g links.
	Next *NavLink `json:"next,omitempty"`
	// Prev is for setting up navigation in hugo templates, e.g. links.
	Prev *NavLink `json:"prev,omitempty"`
}

// Section is a generic part of the book
type Section struct {
	// ID is used to uniquely referenece a section. Can be suffixed
	// with an index to reference a specific text paragraph.
	ID string `json:"id"`
	// Kind is is the kind of section
	Kind string `json:"kind"`
	// Frag is the url fragment that references this segment, sans the #
	Frag string `json:"frag"`
	// Title is used for  section headings
	Title string `json:"title"`
	// Link is the direct link to the content page for this section (if it exists)
	Link string `json:"link,omitempty"`
	// Text is a list of paragraphs that may contain embedded HTML
	Text []string `json:"text"`
	// Sections are child sections, rendered after the above text.
	Sections []*Section `json:"sections"`

	// Hugo stuff

	// Layout is used to pick the right template file, should match Kind.
	Layout string `json:"layout,omitempty"`
	// Book is for setting up navigation in hugo templates.
	Book int `json:"book,omitempty"`
	// Next is for setting up navigation in hugo templates
	Next *NavLink `json:"next,omitempty"`
	// Prev is for setting up navigation in hugo templates.
	Prev *NavLink `json:"prev,omitempty"`
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
	case "cn":
		s = b.parseSimple(d2, short, "common-notion")
	case "prop", "prop1", "prop2", "prop3":
		s = b.parseProps(d2, short)
	default:
		return fmt.Errorf("elem.%d: invalid d2.N: %q (Def|Post|CN|Prop)", b.Number, d2.N)
	}

	b.Sections = append(b.Sections, &s)
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
		Frag: short + "s",
		Title: cleanContent(d2.Heads[0]),
		Sections: make([]*Section, len(d2.Divs)),
	}

	for _, p := range d2.Paras {
		s.Text = append(s.Text, cleanContent(p))
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
		s.Sections[i] = &Section{
			ID: d3.ID,
			Kind: kind,
			Frag: fmt.Sprintf("%s.%s", short, d3.N),
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
		ID: fmt.Sprintf("elem.%d.prop", b.Number), // TODO wrong for Book X
		Kind: "list:proposition",
		Frag: fragify(short),
		Title: cleanContent(d2.Heads[0]),
		Sections: make([]*Section, len(d2.Divs)),
	}
	for i, d3 := range d2.Divs {
		// XXX Book II also uses type="proposition"
		if d3.Type != "number" && d3.Type != "proposition" {
			log.Fatalf("invalid d3.type: %q (number:proposition)", d3.Type)
		}
		ss := Section{
			ID: d3.ID,
			Frag: "prop." + d3.N,
			Kind: "proposition",
			Title: fmt.Sprintf("Proposition %s.", d3.N),
		}
		for _, d4 := range d3.Divs {
			switch d4.Type {
			case "Enunc":
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the theorem, not %d", len(d4.Paras))
				}
				ss.Sections = append(ss.Sections, &Section{
					ID: d3.ID + ".theorem", // TODO Rationalize ID strat better
					// TODO Fragify theorem
					Kind: "theorem",
					Text: []string{ cleanContent(d4.Paras[0]) },
				})
			case "Proof":
				if len(d4.Paras) < 2 {
					log.Fatalf("Expected some steps for the proof, not %d", len(d4.Paras))
				}
				sss := Section{
					ID: d3.ID + ".proof", // TODO Rationalize ID strat better
					// TODO Fragify proof
					Kind: "proof",
					Text: make([]string, len(d4.Paras)),
				}
				for j, p := range d4.Paras {
					sss.Text[j] = cleanContent(p)
				}
				ss.Sections = append(ss.Sections, &sss)
			case "QED":
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the QED, not %d", len(d4.Paras))
				}
				ss.Sections = append(ss.Sections, &Section{
					Kind: "qed",
					// TODO Fragify QED
					Text: []string{cleanContent(d4.Paras[0])},
				})
			case "porism": // TODO
			case "lemma": // TODO
			default:
				log.Fatalf("invalid d4.type: %q (Enunc|Proof|QED|porism|lemma)", d4.Type)
			}
		}

		s.Sections[i] = &ss
	}
	return s
}

// fragify maps a section "short code" to a url fragment, sans the `#`.
func fragify(short string) string {
	switch short {
	case "def", "post", "cn", "prop":
		return short + "s"
	case "def1", "def2", "def3", "prop1", "prop2", "prop3":
		return short
	}
	log.Fatalf("hash: invalid short-code: %q", short)
	return ""
}
//
//// pathify maps a content identifier [e.g elem.#(.short)?.num(.ext)?]
//// For mapping references by id to the urls on the site
//func pathify(id string) string {
//	log.Fatalf("pathify: TODO: more thought required")
//	return ""
//}

// cleanContent does any content transformation that may be necessary on a node
// in the source text. The one kludge for now is resolving `href` attributes in
// anchor tags to the right path on the site.
func cleanContent(p Node) string {
	s := string(p.Content)
	for _, repl := range repls {
	  s = repl.re.ReplaceAllString(s, repl.target)
	}
	return s
}


var repls = []struct{
  target string
  re *regexp.Regexp
} {
	// Drop the middle `.` for the common notions
	// Convert the common notions from elem.1.c.n.Y => /books/X/#cn.Y
	{ `<a href="/books/1/#cn.${1}">`, regexp.MustCompile(`<a href="#elem\.1\.c\.n\.([1-5])">`)},
	// Convert the props from elem.X.Y => /books/X/#prop.Y
	// TODO Figure out the sub-lemmas, etc.
	{ `<a href="/books/${1}/#prop.${2}">`, regexp.MustCompile(`<a href="#elem\.(\d+)\.([0-9]+)">`)},
	// Convert the rest from elem.X.* => /books/X/#*
	{ `<a href="/books/${1}/#${2}">`, regexp.MustCompile(`<a href="#elem\.(\d+)\.([a-z0-9.]+)">`)},
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
