package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ericchiang/css"

	"golang.org/x/net/html"
)

type PartOfSpeech string

const (
	SoftHyphen = "­"

	POSNoun      PartOfSpeech = "Substantiv"
	POSAdjective PartOfSpeech = "Adjektiv"
	POSVerb      PartOfSpeech = "Verb"

	KeywordGrammar        = "Grammatik"
	KeywordNominative     = "Nominativ"
	KeywordPastParticiple = "Partizip II"
)

var (
	ErrTooManyNodes        = errors.New("Found too many nodes")
	ErrMissingNode         = errors.New("Could not find node")
	ErrUnknownPartOfSpeech = errors.New("Unknown part of speech")
	ErrNoGrammar           = errors.New("Could not find grammar")
	ErrNoAdjectiveForms    = errors.New("Could not extract adjective forms")
	ErrNoNounForms         = errors.New("Could not extract noun forms")
	ErrNoVerbForms         = errors.New("Could not extract verb forms")
	ErrUnknownAuxilliary   = errors.New("Could not find auxilliary verb")

	BaseURL = "http://www.duden.de/rechtschreibung"
	//BaseURL = "http://localhost:8000"

	// adjective
	RegexpAdjective     = regexp.MustCompile(`Steigerungsformen:\s*(.*)\s*,\s*(.*)\s*`)
	RegexpThirdSingular = regexp.MustCompile(`er/sie/es (.*)`)
	RegexpAuxilliary    = regexp.MustCompile(`»(.*)«`)
)

func main() {
	for _, arg := range os.Args[1:] {
		word, _ := NewWord(arg)
		fmt.Println(word.BaseWord())
		fmt.Println(word.PartOfSpeech())
		fmt.Println(word.ExtractForms())
		fmt.Println()
	}
}

type Word struct {
	root *html.Node
}

func NewWord(word string) (w Word, err error) {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := c.Get(fmt.Sprintf("%s/%s", BaseURL, word))
	if err != nil {
		return
	}

	w.root, err = html.Parse(res.Body)
	return
}

func (w *Word) BaseWord() (string, error) {
	sel, err := css.Compile("h1")
	if err != nil {
		return "", err
	}

	n, err := exactlyOne(sel, w.root)
	if err != nil {
		return "", err
	}

	return text(n), nil
}

func (w *Word) PartOfSpeech() (PartOfSpeech, error) {
	sel, err := css.Compile("strong.lexem")
	if err != nil {
		return "", err
	}

	for _, n := range sel.Select(w.root) {
		s := text(n)
		switch true {
		case strings.Contains(s, string(POSNoun)):
			return POSNoun, nil
		case strings.Contains(s, string(POSVerb)):
			return POSVerb, nil
		case strings.Contains(s, string(POSAdjective)):
			return POSAdjective, nil
		}
	}

	return "", ErrUnknownPartOfSpeech
}

func (w *Word) ExtractForms() ([]string, error) {
	pos, err := w.PartOfSpeech()
	if err != nil {
		return nil, err
	}

	g, err := grammarNode(w.root)
	if err != nil {
		return nil, err
	}

	switch pos {
	case POSAdjective:
		return extractAdjectiveForms(g)
	case POSNoun:
		return extractNounForms(g)
	case POSVerb:
		base, err := w.BaseWord()
		if err != nil {
			return nil, err
		}
		forms, err := extractVerbForms(g)
		if err != nil {
			return nil, err
		}
		return append([]string{base}, forms...), nil
	}

	return nil, ErrNoGrammar
}

func extractAdjectiveForms(g *html.Node) ([]string, error) {
	sel, err := css.Compile(".lexem")
	if err != nil {
		return nil, err
	}

	n, err := exactlyOne(sel, g)
	if err != nil {
		return nil, err
	}

	matches := RegexpAdjective.FindStringSubmatch(text(n))
	if matches == nil {
		return nil, ErrNoAdjectiveForms
	}

	return matches[1:], nil
}

func extractNounForms(g *html.Node) ([]string, error) {
	sel, err := css.Compile("tbody tr")
	if err != nil {
		return nil, err
	}

	for _, n := range sel.Select(g) {
		c := n.FirstChild
		if strings.Contains(text(c), KeywordNominative) {
			var out []string
			for c = c.NextSibling; c != nil; c = c.NextSibling {
				out = append(out, text(c))
			}
			return out, nil
		}
	}

	return nil, ErrNoNounForms
}

func extractVerbForms(g *html.Node) ([]string, error) {
	aux, err := extractAuxilliary(g)
	if err != nil {
		return nil, err
	}

	sel, err := css.Compile("tbody")
	if err != nil {
		return nil, err
	}

	n := sel.Select(g)
	if len(n) != 3 {
		return nil, ErrNoVerbForms
	}

	out := make([]string, 3)

	// present tense
	for r := n[0].FirstChild; r != nil; r = r.NextSibling {
		s := text(r.FirstChild.NextSibling)
		if matches := RegexpThirdSingular.FindStringSubmatch(s); matches != nil {
			out[0] = matches[1]
		}
	}

	// past tense
	for r := n[1].FirstChild; r != nil; r = r.NextSibling {
		s := text(r.FirstChild.NextSibling)
		if matches := RegexpThirdSingular.FindStringSubmatch(s); matches != nil {
			out[1] = matches[1]
		}
	}

	// perfect tense
	for r := n[2].FirstChild; r != nil; r = r.NextSibling {
		if strings.Contains(text(r.FirstChild), KeywordPastParticiple) {
			out[2] = fmt.Sprintf("%s %s", text(r.LastChild), aux)
		}
	}

	for i := range out {
		if out[i] == "" {
			return out, ErrNoVerbForms
		}
		out[i] = strings.TrimSpace(out[i])
	}

	return out, nil
}

func extractAuxilliary(g *html.Node) (string, error) {
	sel, err := css.Compile(".lexem")
	if err != nil {
		return "", err
	}

	n, err := exactlyOne(sel, g)
	if err != nil {
		return "", ErrUnknownAuxilliary
	}

	matches := RegexpAuxilliary.FindStringSubmatch(text(n))
	if matches == nil {
		return "", ErrUnknownAuxilliary
	}

	return normalizeAuxilliary(matches[1])
}

func normalizeAuxilliary(word string) (string, error) {
	if word == "hat" {
		return "haben", nil
	}
	if word == "ist" {
		return "sein", nil
	}
	return "", ErrUnknownAuxilliary
}

func text(n *html.Node) string {
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			b.WriteString(strings.Replace(c.Data, SoftHyphen, "", -1))
		default:
			b.WriteString(text(c))
		}
	}
	return b.String()
}

func grammarNode(root *html.Node) (*html.Node, error) {
	sel, err := css.Compile("h2")
	if err != nil {
		return nil, err
	}

	for _, n := range sel.Select(root) {
		if strings.Contains(text(n), KeywordGrammar) {
			// search up for the nearest <section>
			for !(n == nil || (n.Data == "section" && n.Type == html.ElementNode)) {
				n = n.Parent
			}
			return n, nil
		}
	}

	return nil, ErrNoGrammar
}

func exactlyOne(sel *css.Selector, n *html.Node) (*html.Node, error) {
	ns := sel.Select(n)
	if len(ns) > 1 {
		return nil, ErrTooManyNodes
	} else if len(ns) == 0 {
		return nil, ErrMissingNode
	}
	return ns[0], nil
}
