// euclid extracts content from the Heath version
//
// TODO:
//	Differentiating Propositions from Proof steps beyond book 1
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

	books := make([]*Book, 0)
	for _, arg := range flag.Args() {
		debug("reading: %s", arg)
		buf, err := ioutil.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}

		// Decode the raw XML structure
		var vol Volume
		err = xml.Unmarshal(buf, &vol)
		if err != nil {
			log.Fatal(err)
		}

		// Transfrom to the more strucuted Book type.
		for _, d1 := range vol.Text.Body.Divs {
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
		path := fmt.Sprintf(filepath.Join(*dir, "book%02d.json"), b.Num)
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

type Book struct {
	Title string `json:"title"`
	Num int `json:"num"`
	Sections []Section `json:"sections"`

	// TODO Remove these now??
	Definitions []Definition `json:"-"`
	Postulates []Postulate `json:"-"`
	CommonNotions []CommonNotion `json:"-"`
	Propositions []Proposition `json:"-"`
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
	b.Num = n
	b.Title = fmt.Sprintf("Book %s", roman(n))

	for _, d2 := range d1.Divs {
		b.parseSection(d2)
	}
	return nil
}

// parseSection TODO
func (b *Book) parseSection(d2 Div2) error {
	debug("  %#v\n", d2.div)
	if d2.Type != "type" {
		return fmt.Errorf("invalid d2.type: %q (book)", d2.Type)
	}

	switch d2.Type {
	case "type":
		switch d2.N {
			// Book X interleaves sections as  `Def #` or `Prop #`
			case "Def":
				s, defs := parseDefs(d2)
				s.ID = fmt.Sprintf("elem.%d.def", b.Num)
				b.Sections = append(b.Sections, s)
				for _, d := range defs {
					b.Definitions = append(b.Definitions, d)
				}
			case "Def 1", "Def 2", "Def 3":
				s, defs := parseDefs(d2)
				// TODO Prop def ID for book X
				b.Sections = append(b.Sections, s)
				for _, d := range defs {
					b.Definitions = append(b.Definitions, d)
				}
			case "Post":
				s, posts := parsePosts(d2)
				s.ID = fmt.Sprintf("elem.%d.post", b.Num)
				b.Sections = append(b.Sections, s)
				for _, p := range posts {
					b.Postulates = append(b.Postulates, p)
				}
			case "CN":
				s, cns := parseCNs(d2)
				s.ID = fmt.Sprintf("elem.%d.c.n", b.Num) // TODO Don't like the c.n.
				b.Sections = append(b.Sections, s)
				for _, cn := range cns {
					b.CommonNotions = append(b.CommonNotions, cn)
				}
			case "Prop":
				s, props := parseProps(d2)
				s.ID = fmt.Sprintf("elem.%d.prop", b.Num)
				b.Sections = append(b.Sections, s)
				for _, p := range props {
					b.Propositions = append(b.Propositions, p)
				}
			case "Prop 1", "Prop 2", "Prop 3":
				s, props := parseProps(d2)
				// TODO Handle the prop IDs for book X
				b.Sections = append(b.Sections, s)
				for _, p := range props {
					b.Propositions = append(b.Propositions, p)
				}
			default:
				return fmt.Errorf("invalid d2.N: %q (Def|Post|CN|Prop)", d2.N)
		}
	default:
		return fmt.Errorf("invalid type: %q (type)", d2.Type)
	}
	return nil
}

// Definition TODO
type Definition struct {
	ID string `json:"id"`
	Text string `json:"text"`
}

func parseDefs(d2 Div2) (Section, []Definition) {
	s := Section{
		// TODO ID based on the book.num
		Kind: "list:definition",
		Title: "Definitions", // TODO d2.Head?
		Sections: make([]Section, len(d2.Divs)),
	}

	a := make([]Definition, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:definition)", d3.Type)
		}
		if len(d3.Paras) != 1 {
			// TODO V.Def.17 has two paragraphs
			fmt.Fprintf(os.Stderr, "warn: %s: wrong # of d3.paras: %d (1:definition)\n", d3.ID, len(d3.Paras))
		}
		// TODO Need to check for <terms> in the list
		a[i] = Definition{d3.ID, cleanPara(d3.Paras[0])}
		s.Sections[i] = Section{
			ID: d3.ID,
			Kind: "definition",
			Title: fmt.Sprintf("Definition %d", i+1), // TODO: d3.head?
			Text: []string{a[i].Text},
		}
		debug("d3:%s %s", d3.ID, a[i].Text)
	}
	return s, a
}

// Postulate TODO
type Postulate struct {
	ID string `json:"id"`
	Text string `json:"text"`
}

func parsePosts(d2 Div2) (Section, []Postulate) {
	s := Section{
		// TODO ID
		Kind: "list:postulate",
		Title: "Postulates",
		Sections: make([]Section, len(d2.Divs)),
	}
	a := make([]Postulate, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("%s: invalid d3.type: %q (number:postulate)", d3.ID, d3.Type)
		}
		content := cleanPara(d3.Paras[0])
		if len(d3.Paras) != 1 {
			// XXX: Book 1, Postulate 1 starts w/ "Let the following be postulated:"
			if i != 0 {
				log.Fatalf("%s: wrong # of d3.paras: %d (1:postulate)", d3.ID, len(d3.Paras))
			}
			// Promote the starter text outside the first postulate
			s.Text = []string{content}
			content = cleanPara(d3.Paras[1])
		}
		a[i] = Postulate{d3.ID, content}
		s.Sections[i] = Section{
			ID: d3.ID,
			Kind: "postulate",
			Title: fmt.Sprintf("Postulate %d", i+1),
			Text: []string{content},
		}
		debug("d3:%s: %s", d3.ID, a[i].Text)
	}
	return s, a
}

// CommonNotion TODO
type CommonNotion struct {
	ID string `json:"id"`
	Text string `json:"text"`
}
func parseCNs(d2 Div2) (Section, []CommonNotion) {
	s := Section{
		// TODO ID
		Kind: "list:common-notion",
		Title: "Common Notions",
		Sections: make([]Section, len(d2.Divs)),
	}
	a := make([]CommonNotion, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:common-notion)", d3.Type)
		}
		if len(d3.Paras) != 1 {
			log.Fatalf("%s: wrong # of d3.paras: %d (1:common-notion)", d3.ID, len(d3.Paras))
		}
		a[i] = CommonNotion{d3.ID, cleanPara(d3.Paras[0])}
		s.Sections[i] = Section{
			ID: d3.ID,
			Kind: "common-notion",
			Title: fmt.Sprintf("Common Notion %d", i+1),
			Text: []string{a[i].Text},
		}
		debug("d3:%s: %s", d3.ID, a[i].Text)
	}
	return s, a
}

// Proposition TODO
type Proposition struct {
	ID string `json:"id"`
	Claim string `json:"claim,omitempty"`// TODO Enunciation?
	Proof []string `json:"proof,omitempty"`
	QED string `json:"qed,omitempty"`
	// TODO Remove
	Text string `json:"text,omitempty"`
}
func parseProps(d2 Div2) (Section, []Proposition) {
	s := Section{
		// TODO ID
		Kind: "list:proposition",
		Title: "Propositions", // TODO d2.Head because of book X
		Sections: make([]Section, len(d2.Divs)),
	}
	a := make([]Proposition, len(d2.Divs))
	for i, d3 := range d2.Divs {
		// XXX Book II also uses type="proposition"
		if d3.Type != "number" && d3.Type != "proposition" {
			log.Fatalf("invalid d3.type: %q (number:proposition)", d3.Type)
		}

		prop := Proposition{ID: d3.ID}
		ss := Section{
			ID: d3.ID,
			Kind: "proposition",
			Title: fmt.Sprintf("Proposition %d", i+1),
		}
		for _, d4 := range d3.Divs {
			switch d4.Type {
			case "Enunc":
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the claim, not %d", len(d4.Paras))
				}
				prop.Claim = cleanPara(d4.Paras[0])
				ss.Sections = append(ss.Sections, Section {
					// TODO ID
					Kind: "claim",
					Text: []string{ prop.Claim },
				})
			case "Proof":
				if len(d4.Paras) < 2 {
					log.Fatalf("Expected some steps for the proof, not %d", len(d4.Paras))
				}
				prop.Proof = make([]string, len(d4.Paras))
				sss := Section{
					// TODO ID
					Kind: "proof",
					Text: make([]string, len(d4.Paras)),
				}
				for j, p := range d4.Paras {
					prop.Proof[j] = cleanPara(p)
					sss.Text[j] = prop.Proof[j]
				}
				ss.Sections = append(ss.Sections, sss)
			case "QED": // skip
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the QED, not %d", len(d4.Paras))
				}
				prop.QED = cleanPara(d4.Paras[0])
				ss.Sections = append(ss.Sections, Section{
					// TODO ID
					Kind: "qed",
					Text: []string{prop.QED},
				})
			case "porism": // TODO Definitely a new section
			case "lemma": // TODO Definitely a new section
			default:
				log.Fatalf("invalid d4.type: %q (Enunc|Proof|QED|porism|lemma)", d4.Type)
			}
		}

		if len(prop.Proof) == 0 {
			// TODO Can we assume d3.Nodes[0] as Claim and d3.Nodes[1:] as Proof?
			prop.Text = string(d3.Content)
			debug("d3:%s:raw %s", d3.ID, prop.Text)
		} else {
			debug("d3:%s:claim %s\n", d3.ID, prop.Claim)
			debug("d3:%s:proof %s\n", d3.ID, prop.Proof)
		}
		a[i] = prop
		s.Sections[i] = ss
	}
	return s, a
}


// cleanPara is a bit of a regex kludge. It exists to transform embedded
// tags in paragraphs from the source XML to something that is HTML-friendly.
func cleanPara(p Node) string {
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
	{ `<${1}dfn>`, regexp.MustCompile(`<(/?)term>`)},
	{ `<${1}var>`, regexp.MustCompile(`<(/?)emph>`)},
	// TODO Make these into superscripts of the preceding statement
	{ `<a href="#$1">`, regexp.MustCompile(`<ref target="([a-z1-9.]+)" targOrder="U">`) },
	{ `</a>`, regexp.MustCompile(`</ref>`) },
	{ ``, regexp.MustCompile(`<figure />`) },
	{ ``, regexp.MustCompile(`<(/?)hi( rend="center")?>`) },
	{ ``, regexp.MustCompile(`<(/?)hi( rend="center")?>`) },
	{ ``, regexp.MustCompile(`<[p|l]b n="\d+" />`) },
	// TODO anything with <note>'s?
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
