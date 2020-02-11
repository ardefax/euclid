// euclid extracts content from the Heath version
//
// TODO:
//	Differentiating Propositions from Proof steps beyond book 1
//	HTML-transforms
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

// Book TODO Later books (e.g. X) interleave definitions and propositions
type Book struct {
	Num int `json:"num"`

	Definitions []Definition `json:"definitions"`
	Postulates []Postulate `json:"postulates"`
	CommonNotions []CommonNotion `json:"common_notions"`
	Propositions []Proposition `json:"propositions"`
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

	for _, d2 := range d1.Divs {
		debug("  %#v\n", d2.div)
		if d2.Type != "type" {
			return fmt.Errorf("invalid d2.type: %q (book)", d1.Type)
		}

		switch d2.Type {
		case "type":
			switch d2.N {
				// TODO Book X interleaves definitions and propostions
				// in sections as `Def #` or `Prop #`
				case "Def":
					b.Definitions = defs(d2)
				case "Post":
					b.Postulates = posts(d2)
				case "CN":
					b.CommonNotions = cns(d2)
				case "Prop":
					b.Propositions = props(d2)
				default:
					return fmt.Errorf("invalid d2.N: %q (Def|Post|CN|Prop)", d2.N)
			}
		default:
			return fmt.Errorf("invalid type: %q (type)", d2.Type)
		}
	}
	return nil
}

// Definition TODO
type Definition struct {
	ID string `json:"id"`
	Text string `json:"text"`
}

func defs(d2 Div2) []Definition {
	a := make([]Definition, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:definition)", d3.Type)
		}
		if len(d3.Paras) != 1 {
			log.Fatalf("%s: wrong # of d3.paras: %d (1:definition)", d3.ID, len(d3.Paras))
		}
		// TODO Need to check for <terms> in the list
		a[i] = Definition{d3.ID, cleanPara(d3.Paras[0])}
		debug("d3:%s %s", d3.ID, a[i].Text)
	}
	return a
}

// Postulate TODO
type Postulate struct {
	ID string `json:"id"`
	Text string `json:"text"`
}

func posts(d2 Div2) []Postulate {
	a := make([]Postulate, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("%s: invalid d3.type: %q (number:postulate)", d3.ID, d3.Type)
		}
		content := cleanPara(d3.Paras[0])
		if len(d3.Paras) != 1 {
			// XXX: Book 1, Postulate 1 starts w/ "Let the following be postulated:"
			// This was manually injected as part of the page layout for now :/
			//log.Fatalf("%s: wrong # of d3.paras: %d (1:postulate)", d3.ID, len(d3.Paras))
			fmt.Fprintf(os.Stderr, "warn: %s: wrong # of d3.paras: %d (1:postulate)\n", d3.ID, len(d3.Paras))
			content = cleanPara(d3.Paras[1])
		}
		a[i] = Postulate{d3.ID, content}
		debug("d3:%s: %s", d3.ID, a[i].Text)
	}
	return a
}

// CommonNotion TODO
type CommonNotion struct {
	ID string `json:"id"`
	Text string `json:"text"`
}
func cns(d2 Div2) []CommonNotion {
	a := make([]CommonNotion, len(d2.Divs))
	for i, d3 := range d2.Divs {
		if d3.Type != "number" {
			log.Fatalf("invalid d3.type: %q (number:common-notion)", d3.Type)
		}
		if len(d3.Paras) != 1 {
			log.Fatalf("%s: wrong # of d3.paras: %d (1:common-notion)", d3.ID, len(d3.Paras))
		}
		a[i] = CommonNotion{d3.ID, cleanPara(d3.Paras[0])}
		debug("d3:%s: %s", d3.ID, a[i].Text)
	}
	return a
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
func props(d2 Div2) []Proposition {
	a := make([]Proposition, len(d2.Divs))
	for i, d3 := range d2.Divs {
		// TODO Book II also uses type="proposition"
		if d3.Type != "number" && d3.Type != "proposition" {
			log.Fatalf("invalid d3.type: %q (number:proposition)", d3.Type)
		}

		prop := Proposition{ID: d3.ID}
		for _, d4 := range d3.Divs {
			switch d4.Type {
			case "Enunc":
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the claim, not %d", len(d4.Paras))
				}
				prop.Claim = cleanPara(d4.Paras[0])
			case "Proof":
				if len(d4.Paras) < 2 {
					log.Fatalf("Expected some steps for the proof, not %d", len(d4.Paras))
				}
				prop.Proof = make([]string, len(d4.Paras))
				for j, p := range d4.Paras {
					prop.Proof[j] = cleanPara(p)
				}
			case "QED": // skip
				if len(d4.Paras) != 1 {
					log.Fatalf("Expected 1 paragraph for the QED, not %d", len(d4.Paras))
				}
				prop.QED = cleanPara(d4.Paras[0])
			case "porism": // TODO
			case "lemma": // TODO
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
	}
	return a
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
